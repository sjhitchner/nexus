package publish

import (
	"fmt"
	"sync"
)

type LogPublisher struct {
	sync.RWMutex
}

func (t LogPublisher) Publish(b []byte) {
	t.Lock()
	defer t.Unlock()

	fmt.Println("====")
	fmt.Print(string(b))
	fmt.Println("====")
}
