package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	. "github.com/sjhitchner/infosphere/common"
	"log"
	"net/http"
	"sync"
	"time"
)

func Client1() {
	origin := "http://localhost/"
	url := "ws://localhost:12345/websocket/1"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int

	for {
		if n, err = ws.Read(msg); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received: %s.\n", msg[:n])
	}
}

func Client2() {
	origin := "http://localhost/"
	url := "ws://localhost:12345/websocket/2"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int

	for {
		if n, err = ws.Read(msg); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received: %s.\n", msg[:n])
	}
}

func main() {

	go func() {
		select {
		case <-time.After(5 * time.Second):
			go Client1()
			go Client2()
		}
	}()

	channel1 := make(chan Payload, 10)
	go func() {
		for {
			payload := Payload{
				Bucket:    "test1",
				Data:      "data1",
				CreatedAt: time.Now(),
			}
			channel1 <- payload
		}
	}()
	channel2 := make(chan Payload, 10)
	go func() {
		for {
			payload := Payload{
				Bucket:    "test2",
				Data:      "data2",
				CreatedAt: time.Now(),
			}
			channel2 <- payload
		}
	}()

	wm := NewWebsocketManager()
	wm.AddPath("1", channel1)
	wm.AddPath("2", channel2)
	http.Handle("/websocket/", websocket.Handler(wm.Handler))
	//http.Handle("/websocket/", websocket.Handler(TestHandler(channel)))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func TestHandler(channel chan Payload) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		for payload := range channel {
			ws.Write(payload.Bytes())
		}
	}
}

type WebsocketManager struct {
	sync.RWMutex
	channels map[string]chan Payload
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		sync.RWMutex{},
		make(map[string]chan Payload),
	}
}

func (t *WebsocketManager) AddPath(path string, channel chan Payload) error {
	t.Lock()
	defer t.Unlock()

	fullPath := fullPath(path)
	if _, ok := t.channels[fullPath]; ok {
		return fmt.Errorf("Path [%s] already in use", fullPath)
	}

	log.Printf("Added [%s]\n", fullPath)
	t.channels[fullPath] = channel
	return nil
}

func (t *WebsocketManager) RemovePath(path string) error {
	t.Lock()
	defer t.Unlock()

	fullPath := fullPath(path)
	if _, ok := t.channels[fullPath]; !ok {
		return fmt.Errorf("Path [%s] not configured", fullPath)
	}

	close(t.channels[fullPath])
	delete(t.channels, fullPath)

	return nil
}

func (t *WebsocketManager) Handler(ws *websocket.Conn) {
	path := ws.Request().URL.Path

	log.Printf("Handling [%s]\n", path)

	t.RLock()
	channel, ok := t.channels[path]
	t.RUnlock()

	if !ok {
		err := fmt.Errorf("Invalid path [%s]", path)
		if _, err := ws.Write([]byte(err.Error())); err != nil {
			log.Printf("Error writing to ws %v", err)
		}
	}

	log.Printf("Starting connection to %s...", path)
	defer log.Printf("Stopping connection to %s...", path)
	for payload := range channel {
		ws.Write(payload.Bytes())
	}
}

func fullPath(path string) string {
	return fmt.Sprintf("/websocket/%s", path)
}
