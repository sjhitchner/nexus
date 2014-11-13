package aggregator

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

const (
	DELIMITER = byte(10)
)

type Publisher interface {
	Publish(b []byte)
}

type Aggregator interface {
	Sink(payload interface{})
	Start()
	Shutdown()
}

type aggregator struct {
	wg         sync.WaitGroup
	channel    chan interface{}
	publisher  Publisher
	bufferSize int
	numWorkers int
	timeout    time.Duration
}

func NewAggregator(bufferSize int, timeout time.Duration, numWorkers int, publisher Publisher) *aggregator {
	return &aggregator{
		wg:         sync.WaitGroup{},
		channel:    make(chan interface{}, 1000),
		publisher:  publisher,
		bufferSize: bufferSize,
		numWorkers: numWorkers,
		timeout:    timeout,
	}
}

func (t *aggregator) Start() {
	for i := 0; i < t.numWorkers; i++ {
		go t.worker()
	}
}

func (t *aggregator) Shutdown() {
	close(t.channel)
	t.wg.Wait()
	log.Println("Aggregator Shutdown")
}

func (t *aggregator) Sink(data interface{}) {
	t.channel <- data
}

func (t *aggregator) worker() {
	t.wg.Add(1)
	defer t.wg.Done()

	buffer := make([]byte, t.bufferSize)
	counter := 0

	timeoutChannel := time.After(t.timeout)
	for {
		select {
		case payload, ok := <-t.channel:
			if !ok {
				log.Println("Aggregator: channel closed")
				if counter > 0 {
					t.publisher.Publish(buffer[:counter])
				}
				break
			}

			log.Println("Aggregator: received message")
			b, err := json.Marshal(payload)
			if err != nil {
				log.Println("unable to jsonify payload")
				continue
			}
			payloadSize := len(b) + 1

			if counter+payloadSize >= t.bufferSize {
				log.Printf("Fill rate %d/%d=%.02f\n", counter, t.bufferSize, float32(counter)/float32(t.bufferSize))
				t.publisher.Publish(buffer[:counter])
				timeoutChannel = time.After(t.timeout)
				counter = 0
			}

			for i := 0; i < payloadSize-1; i++ {
				buffer[counter] = b[i]
				counter++
			}
			buffer[counter] = DELIMITER
			counter++

		case <-timeoutChannel:
			log.Println("Aggregator: timeout")
			if counter > 0 {
				t.publisher.Publish(buffer[:counter])
				timeoutChannel = time.After(t.timeout)
				counter = 0
			}
		}
	}
}

/*


		if t.publisher != nil {
			t.publisher.Publish(b)
		}


	if len(t.channel) == t.batchSize {
		close(t.channel)

		t.aggregate(t.channel)
		t.channel = make(chan Payload, t.batchSize)
	}
}

func (t *aggregator) aggregate(channel chan Payload) {
	t.wg.Add(1)
	defer t.wg.Done()


	for payload := range channel {

		if t.publisher != nil {
			t.publisher.Publish(b)
		}
	}
}
*/
