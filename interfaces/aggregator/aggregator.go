package aggregator

import (
	"encoding/json"
	. "github.com/sjhitchner/infosphere/domain"
	"log"
	"sync"
)

const (
	//BUFFER_SIZE = 128000
	BUFFER_SIZE = 256
	DELIMITER   = byte(10)
)

type Aggregator interface {
	Sink(data Payload)
	Start()
	Shutdown()
}

type aggregator struct {
	wg        sync.WaitGroup
	channel   chan Payload
	publisher Publisher
}

func NewAggregator(publisher Publisher) *aggregator {
	return &aggregator{
		wg:        sync.WaitGroup{},
		channel:   make(chan Payload, 1000),
		publisher: publisher,
	}
}

func (t *aggregator) Start() {
	go t.worker()
}

func (t *aggregator) Shutdown() {
	close(t.channel)
	t.wg.Wait()
	log.Println("Aggregator Shutdown")
}

func (t *aggregator) Sink(data Payload) {
	t.channel <- data
}

func (t *aggregator) worker() {
	t.wg.Add(1)
	defer t.wg.Done()

	buffer := make([]byte, BUFFER_SIZE)
	counter := 0

	for payload := range t.channel {
		b, err := json.Marshal(payload)
		if err != nil {
			log.Println("unable to jsonify payload")
			continue
		}
		payloadSize := len(b) + 1

		if counter+payloadSize >= BUFFER_SIZE {
			log.Printf("Fill rate %d/%d=%.02f\n", counter, BUFFER_SIZE, float32(counter)/float32(BUFFER_SIZE))
			go t.publisher.Publish(buffer)
			buffer = make([]byte, BUFFER_SIZE)
			counter = 0
		}

		for i := 0; i < payloadSize-1; i++ {
			buffer[counter] = b[i]
			counter++
		}
		buffer[counter] = DELIMITER
		counter++
	}

	if counter > 0 {
		t.publisher.Publish(buffer)
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
