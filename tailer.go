package trash

import (
	"io"
	"log"
	"time"
)

type TailFilterFunc func(s Span) bool

type Tailer interface {
	Write(spans Spans) error
	io.Closer
}

type clockTick struct {
	m map[string][]Span
}

type clockTicks []*clockTick // represent 60 miunte

type MemTailer struct {
	filterFunc TailFilterFunc
	ticks      clockTicks
	invertMap  map[string]*clockTick
}

func NewMemTailer(flushDuration time.Duration, filterFunc TailFilterFunc) *MemTailer {
	ticks := make([]*clockTick, 0, 60)
	m := &MemTailer{
		filterFunc: filterFunc,
		ticks:      ticks,
	}
	t := time.NewTicker(flushDuration)
	doFlush := func(tick map[string][]Span) error {
		return nil
	}

	for range t.C {
		idx := time.Now().Unix() / int64(flushDuration.Seconds()) % int64(len(m.ticks))
		tick := ticks[idx]
		if err := doFlush(tick.m); err != nil {
			log.Print(err)
		}
	}
	return m
}

func (m *MemTailer) Write(spans Spans) error {
	for _, span := range spans {
		m.dispatch(*span)
	}
	return nil
}

func (m *MemTailer) dispatch(span Span) error {
	tid := span.TraceID
	tick := m.invertMap[tid]
	tick.m[tid] = append(tick.m[tid], span)
	return nil
}

func (m *MemTailer) Close() error {
	return nil
}
