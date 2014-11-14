package multiplex

import (
	"github.com/sjhitchner/nexus/interfaces/sink"
	"sync"
)

type Multiplexer interface {
	AddSink(sink ...sink.Sink)
	Multiplex(path string, payload interface{})
}

type multiplexer struct {
	sync.RWMutex
	sinks []sink.Sink
}

func NewMultiplexer() Multiplexer {
	return &multiplexer{
		sync.RWMutex{},
		make([]sink.Sink, 0, 5),
	}
}

func (t *multiplexer) AddSink(sink ...sink.Sink) {
	t.Lock()
	defer t.Unlock()
	t.sinks = append(t.sinks, sink...)
}

func (t *multiplexer) Multiplex(path string, payload interface{}) {
	t.RLock()
	defer t.RUnlock()

	for _, s := range t.sinks {
		s.Sink(path, payload)
	}
}
