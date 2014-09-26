package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	//. "github.com/sjhitchner/infosphere/common"
	"encoding/json"
	"github.com/sjhitchner/infosphere/interfaces/router"
	"log"
	"net/http"
	//"sync"
	"io"
	"time"
)

type Payload struct {
	Bucket    string
	Data      string
	CreatedAt time.Time
}

func (t Payload) String() string {
	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}

func (t Payload) Bytes() []byte {
	b, _ := json.MarshalIndent(t, "", "  ")
	return b
}

var id = 0

func Client1() {
	id++

	origin := "http://localhost/"
	url := fmt.Sprintf("ws://localhost:12345/websocket/%d", id)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int

	done := time.After(3 * time.Second)
	for {
		select {
		case <-done:
			ws.Close()
			return
		default:
			if n, err = ws.Read(msg); err != nil {
				if err == io.EOF {
					ws.Close()
					return
				}
				log.Println(err)
			}
			fmt.Printf("Received: %s.\n", msg[:n])
		}
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
		for {
			select {
			case <-time.After(5 * time.Second):
				go Client1()
				//go Client2()
			}
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				var payload = Payload{
					Bucket:    "test1",
					Data:      "data1",
					CreatedAt: time.Now(),
				}
				router.Route("/websocket/1", payload)
			}
		}
	}()
	/*
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
	*/

	wm := NewWebsocketManager()
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
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{}
}

func (t *WebsocketManager) Handler(ws *websocket.Conn) {
	path := ws.Request().URL.Path

	log.Printf("Handling [%s]\n", path)

	receiver := router.NewChannel(1)
	if err := router.AddChannel(path, receiver); err != nil {
		err := fmt.Errorf("Invalid path [%s]", path)
		log.Printf(err.Error())
		handleError(ws, err)
		return
	}

	log.Printf("Starting connection to %s...", path)
	defer log.Printf("Stopping connection to %s...", path)
	defer func() {
		receiver.Close()
		for _ = range receiver.Dequeue() {
		}
	}()

	/*
		for payload := range receiver.Dequeue() {
			b, err := json.Marshal(payload)
			if err != nil {
				handleError(ws, err)
			}
			if _, err := ws.Write(b); err != nil {
				handleError(ws, err)
			}
		}
	*/
}

func handleError(ws *websocket.Conn, err error) {
	if _, err := ws.Write([]byte(err.Error())); err != nil {
		log.Printf("Error writing error to ws %v", err)
	}
}

func fullPath(path string) string {
	return fmt.Sprintf("/websocket/%s", path)
}
