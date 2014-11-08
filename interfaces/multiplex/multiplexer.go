package multiplex

import (
	. "github.com/sjhitchner/nexus/domain"
	//"log"
	"sync"
)

type Multiplexer interface {
	AddSink(sink ...Sink)
	Multiplex(payload Payload)
}

type multiplexer struct {
	sync.RWMutex
	sinks []Sink
}

func NewMultiplexer() Multiplexer {
	return &multiplexer{
		sync.RWMutex{},
		make([]Sink, 0, 5),
	}
}

func (t *multiplexer) AddSink(sink ...Sink) {
	t.Lock()
	defer t.Unlock()
	t.sinks = append(t.sinks, sink...)
}

func (t *multiplexer) Multiplex(payload Payload) {
	t.RLock()
	defer t.RUnlock()

	for _, sink := range t.sinks {
		sink.Sink(payload)
	}
}
