package aggregator

import (
	"fmt"
	. "github.com/sjhitchner/infosphere/common"
	"sync"
	"testing"
	"time"
)

var plock = sync.RWMutex{}

type TestPublisher struct {
}

func (t TestPublisher) Publish(b []byte) {
	plock.Lock()
	defer plock.Unlock()

	fmt.Println("====")
	fmt.Println(string(b))
	fmt.Println("====")
}

func TestAggregator(t *testing.T) {

	agg := NewAggregator(2)
	agg.AddPublisher(TestPublisher{})

	for i := 0; i < 30; i++ {
		agg.Sink(Payload{"test", "test", time.Now()})
	}

	agg.Shutdown()
}
