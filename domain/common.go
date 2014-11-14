package interfaces

import (
	"encoding/json"
	"time"
)

type Path string

type Filter interface {
	Name() string
	Filter(payload Payload) error
	Buckets() []Path
	Enable()
	Disable()
	IsEnabled() bool
}

type Payload struct {
	Path       Path        `json:"path", codec:"path"`
	Data       interface{} `json:"data", codec:"data"`
	ReceivedAt time.Time   `json:"created_at", codec:"created_at"`
}

func (t Payload) String() string {
	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}

func (t Payload) Bytes() []byte {
	b, _ := json.MarshalIndent(t, "", "  ")
	return b
}
