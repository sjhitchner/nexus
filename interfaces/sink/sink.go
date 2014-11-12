package sink

import (
	"encoding/json"
	"log"
)

type Sink interface {
	Sink(data interface{})
}

type LogSink struct {
}

func (t LogSink) Sink(payload interface{}) {

	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Println("Unable to marshal payload")
		return
	}

	log.Println(string(b))
}
