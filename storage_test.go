package trash

import (
	"github.com/dgraph-io/badger"
	"testing"
)

func TestBadgerStorage_Index(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		t.Fatal(err)
	}
	storage := &BadgerStorage{DB: db}
	spans := []*Span{
		{ID: "foo", TraceID: "prefix1"},
		{ID: "bar", TraceID: "prefix1"},
		{ID: "bar", TraceID: "prefix2"},
	}
	if err := storage.Index(spans); err != nil {
		t.Error(err)
	}

	getspans, err := storage.GetTrace("prefix1")
	if err != nil {
		t.Error(err)
	}
	if len(getspans) != 2 {
		t.Error("expect 2 spans")
	}
}

func TestRotatedBadgerStorage_Index(t *testing.T) {
	storage := &RotatedBadgerStorage{dir: "/tmp", rotation: make(map[int]*BadgerStorage)}
	spans := []*Span{
		{ID: "foo", TraceID: "prefix1"},
		{ID: "bar", TraceID: "prefix1"},
		{ID: "bar", TraceID: "prefix2"},
	}
	if err := storage.Index(spans); err != nil {
		t.Error(err)
	}

	getspans, err := storage.GetTrace("prefix1")
	if err != nil {
		t.Error(err)
	}
	if len(getspans) != 2 {
		t.Error("expect 2 spans")
	}
}
