package interfaces

import (
	"encoding/json"
	"time"
)

type Bucket string

type Filter interface {
	Name() string
	Filter(payload Payload) error
	Buckets() []Bucket
	Enable()
	Disable()
	IsEnabled() bool
}

type Payload struct {
	Bucket    Bucket      `json:"bucket", codec:"bucket"`
	Data      interface{} `json:"data", codec:"data"`
	CreatedAt time.Time   `json:"created_at", codec:"created_at"`
}

func (t Payload) String() string {
	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}

func (t Payload) Bytes() []byte {
	b, _ := json.MarshalIndent(t, "", "  ")
	return b
}

type Publisher interface {
	Publish(b []byte)
}

type Sink interface {
	Sink(data Payload)
}
