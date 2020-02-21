package trash

import (
	"bytes"
	"github.com/dgraph-io/badger"
	"log"
)

type Storage interface {
	Index(spans Spans) error
	GetTrace(id []byte) (Spans, error)
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
		err = wb.Set(b.Bytes(), payload, 0)
		if err != nil {
			log.Print(err)
		}
	}
	return nil
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
