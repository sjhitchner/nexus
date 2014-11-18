package aggregator

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

const (
	DELIMITER = byte(10)
	// DELIMITER = byte(0)
	CHANNEL_DEPTH = 1000
)

type Publisher interface {
	Publish(b []byte)
}

type Aggregator interface {
	Sink(path string, payload interface{})
	Shutdown()
}

type aggregator struct {
	sync.RWMutex
	wg         sync.WaitGroup
	channels   map[string]chan interface{}
	publisher  Publisher
	bufferSize int
	timeout    time.Duration
}

func NewAggregator(bufferSize int, timeout time.Duration, publisher Publisher) *aggregator {
	return &aggregator{
		wg:         sync.WaitGroup{},
		channels:   make(map[string]chan interface{}),
		publisher:  publisher,
		bufferSize: bufferSize,
		timeout:    timeout,
	}
}

// Sink a path and data
func (t *aggregator) Sink(path string, data interface{}) {
	t.RLock()
	channel, ok := t.channels[path]
	t.RUnlock()

	if !ok {
		t.Lock()
		channel = make(chan interface{}, 1000)
		t.Unlock()

		go t.worker(path, channel)
		t.channels[path] = channel
	}
	channel <- data
}

func (t *aggregator) worker(path string, channel chan interface{}) {
	t.wg.Add(1)
	defer t.wg.Done()

	log.Printf("starting worker for [%s]", path)

	buffer := make([]byte, t.bufferSize)
	counter := 0

	timeoutChannel := time.After(t.timeout)
	for {
		log.Println(t.channels)
		select {
		case payload, ok := <-channel:
			if !ok {
				log.Printf("Aggregator: channel closed for [%s]", path)
				if counter > 0 {
					t.publisher.Publish(buffer[:counter])
					counter = 0
				}
				return
			}

			log.Printf("Aggregator: received message for [%s]", path)
			b, err := json.Marshal(payload)
			if err != nil {
				log.Println("unable to jsonify payload")
				continue
			}

			payloadSize := len(b) + 1
			if counter+payloadSize >= t.bufferSize {
				log.Printf("Fill rate %d/%d=%.02f\n", counter, t.bufferSize, float32(counter)/float32(t.bufferSize))
				t.publisher.Publish(buffer[:counter])
				counter = 0
			}

			for i := 0; i < payloadSize-1; i++ {
				buffer[counter] = b[i]
				counter++
			}
			buffer[counter] = DELIMITER
			counter++

			timeoutChannel = time.After(t.timeout)
			log.Println("P:", payload)

		case <-timeoutChannel:
			log.Printf("Aggregator: timeout for [%s]", path)
			t.Lock()
			close(channel)
			delete(t.channels, path)
			t.Unlock()
		}
	}
}

func (t *aggregator) Shutdown() {
	t.Lock()
	defer t.Unlock()

	for path, _ := range t.channels {
		close(t.channels[path])
		delete(t.channels, path)
	}

	t.wg.Wait()
	log.Printf("Shutting down Aggregator %v\n", t.channels)
}

/*


			if !ok {
				log.Printf("Aggregator: channel closed for [%s]", path)
				if counter > 0 {
					t.publisher.Publish(buffer[:counter])
					counter = 0
				}
				return
			}

			log.Printf("Aggregator: received message for [%s]", path)
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
