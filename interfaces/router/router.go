package router

import (
	"container/list"
	"errors"
	"log"
	"sync"
)

type Channel struct {
	sync.RWMutex
	channel chan interface{}
	closed  bool
}

func NewChannel(depth int) *Channel {
	return &Channel{
		sync.RWMutex{},
		make(chan interface{}, depth),
		false,
	}
}

func (t *Channel) Enqueue(data interface{}) error {
	t.RLock()
	defer t.RUnlock()
	if t.closed {
		return errors.New("channel closed")
	}
	t.channel <- data
	return nil
}

func (t *Channel) Dequeue() <-chan interface{} {
	return t.channel
}

func (t *Channel) Close() {
	t.Lock()
	defer t.Unlock()
	t.closed = true
	close(t.channel)
}

type RouteList struct {
	*sync.RWMutex
	*list.List
}

type RouteMap map[string]RouteList

var lock = sync.RWMutex{}
var routes = make(RouteMap)

func AddChannel(path string, channel ...*Channel) error {
	lock.Lock()
	defer lock.Unlock()

	log.Printf("Adding channel for path [%s]\n", path)

	if _, ok := routes[path]; !ok {
		routes[path] = RouteList{
			&sync.RWMutex{},
			list.New(),
		}
	}

	routes[path].Lock()
	defer routes[path].Unlock()
	for _, c := range channel {
		routes[path].PushBack(c)
	}

	return nil
}

func Route(path string, data interface{}) {
	lock.RLock()
	defer lock.RUnlock()

	routes[path].RLock()
	for e := routes[path].Front(); e != nil; e = e.Next() {
		channel := e.Value.(*Channel)
		if err := channel.Enqueue(data); err != nil {
			routes[path].RUnlock()
			routes[path].Lock()
			routes[path].Remove(e)
			routes[path].Unlock()
			routes[path].RLock()
		}
	}
	routes[path].RUnlock()
}
