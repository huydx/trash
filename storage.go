package trash

import (
	"bytes"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
)

type Storage interface {
	Index(spans Spans) error
	GetTrace(id string) (Spans, error)
}

type RotatedBadgerStorage struct {
	dir      string
	rotation map[int]*BadgerStorage // map[yearday]storage
}

func (r RotatedBadgerStorage) Index(spans Spans) error {
	n := time.Now()
	yd := n.YearDay()
	if r.rotation[yd] == nil {
		d := n.Format("20060102")
		db, err := badger.Open(badger.DefaultOptions(filepath.Join(r.dir, d)))
		if err != nil {
			log.Fatal(err)
		}
		r.rotation[yd] = &BadgerStorage{DB: db}
	}
	return r.rotation[yd].Index(spans)
}

func (r RotatedBadgerStorage) GetTrace(id string) (Spans, error) {
	w := sync.WaitGroup{}
	spans := make(Spans, 0)
	mu := sync.Mutex{}
	// fanout to all date partition
	for _, rot := range r.rotation {
		w.Add(1)
		go func() {
			defer w.Done()
			sp, err := rot.GetTrace(id)
			if err != nil {
				log.Print(err)
				return
			}
			mu.Lock()
			spans = append(spans, sp...)
			mu.Unlock()
		}()
	}
	w.Wait()
	return spans, nil
}

type BadgerStorage struct {
	DB *badger.DB
}

func (b BadgerStorage) Index(spans Spans) error {
	wb := b.DB.NewWriteBatch()
	defer wb.Cancel()
	for _, span := range spans {
		b := bytes.NewBuffer([]byte{})
		b.Write([]byte(span.TraceID))
		b.Write([]byte(span.ID))
		payload, err := span.Marshal()
		if err != nil {
			log.Print(err)
			continue
		}
		err = wb.Set(b.Bytes(), payload)
		if err != nil {
			log.Print(err)
		}
	}
	return wb.Flush()
}

func (b BadgerStorage) GetTrace(id string) (Spans, error) {
	prefix := []byte(id)
	spans := make([]*Span, 0)
	b.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			span := Span{}
			err := item.Value(func(v []byte) error {
				return span.Unmarshal(v)
			})
			if err != nil {
				return err
			}
			spans = append(spans, &span)
		}
		return nil
	})
	return spans, nil
}
