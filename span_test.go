package trash

import (
	"testing"
)

func TestSpan_Marshal(t *testing.T) {
	s := &Span{
		ID:      "123",
		TraceID: "123",
	}
	bs, err := s.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) == 0 {
		t.Fatal("expect > 0 bytes")
	}

	unmarshal := &Span{}
	if err = unmarshal.Unmarshal(bs); err != nil {
		t.Fatal(err)
	}
	if unmarshal.ID != "123" {
		t.Fatalf("expect same span id, got %v", unmarshal)
	}
}
