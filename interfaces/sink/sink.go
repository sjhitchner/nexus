package sink

import (
	"log"
)

type LogSink struct {
}

func (t *LogSink) Sink(payload Payload) {
	log.Println(payload)
}
