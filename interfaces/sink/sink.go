package sink

import (
	"encoding/json"
	"log"
)

type Sink interface {
	Sink(path string, data interface{})
}

type LogSink struct {
}

func (t LogSink) Sink(path string, payload interface{}) {

	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Println("Unable to marshal payload")
		return
	}

	log.Printf("LOGSINK %s:%s\n", path, string(b))
}
