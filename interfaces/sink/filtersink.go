package filter

import (
	"fmt"
	"sync"
)

type FilterSink interface {
	Sink(data Payload)
	AddFilter(filter Filter)
	RemoveFilter(filterName string) error
}

type FilterMap map[Bucket][]*Filter

type defaultFilterSink struct {
	sync.RWMutex
	filters FilterMap
}

func NewFilterSink() FilterSink {
	return &defaultFilterSink{
		sync.RWMutex{},
		make(FilterMap),
	}
}

func (t *defaultFilterSink) Sink(data Payload) {
	t.RLock()
	defer t.RUnlock()

	for _, filter := range t.filters[data.Bucket] {
		filter.Filter(payload)
	}
}

func (t *defaultFilterSink) AddFilter(filter Filter) {
	t.Lock()
	defer t.Unlock()

	for _, bucket := range filter.Buckets() {
		_, ok := t.filters[bucket]
		if !ok {
			t.filters[bucket] = make([]*Filter, 0)
		}
		t.filters[bucket] = append(t.filters[bucket], &filter)
	}
}
