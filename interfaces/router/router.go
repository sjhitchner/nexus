package router

import (
	"log"
	"sync"
)

type Connections map[string]chan interface{}
type Routes map[string]Connections

var routes = make(Routes)
var lock = sync.RWMutex{}

func AddRoute(path string, connectionId string, channel chan interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	log.Printf("Adding channel for path [%s]\n", path)

	if _, ok := routes[path]; !ok {
		routes[path] = make(Connections)
	}

	routes[path][connectionId] = channel
	return nil
}

func RemoveRoute(path string, connectionId string) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := routes[path]; !ok {
		return
	}

	channel := routes[path][connectionId]
	close(channel)

	delete(routes[path], connectionId)

	if len(routes[path]) == 0 {
		delete(routes, path)
	}
}

func Route(path string, data interface{}) {
	lock.RLock()
	defer lock.RUnlock()

	route, ok := routes[path]
	if !ok {
		return
	}
	for _, channel := range route {
		channel <- data
	}

	log.Println(routes)
}

/*
	t.RLock()
	defer t.RUnlock()
	select {
	case <-done:
		return
		case <-
	}
	if t.closed {
		return errors.New("channel closed")
	}
	channel <- data
	return nil
}
*/

/*
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

type RouteMap map[string]*Channel

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

	log.Println("Routes", len(routes))
	for _, v := range routes {
		log.Println(v.Len())
	}

	return nil
}

func CloseChannel(path string) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := routes[path]; !ok {
		return
	}

}

func Route(path string, data interface{}) {
	lock.RLock()
	defer lock.RUnlock()

	if _, ok := routes[path]; !ok {
		return
	}

	routes[path].RLock()
	for e := routes[path].Front(); e != nil; e = e.Next() {
		channel := e.Value.(*Channel)
		if err := channel.Enqueue(data); err != nil {
			log.Println("tried to write to closed channel")
		}
	}
	routes[path].RUnlock()
}
*/
