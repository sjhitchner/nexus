package handlers

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"log"
	"net/http"
)

type WebsocketHandler struct {
	channel chan interface{}
}

func NewWebsocketHandler() *WebsocketHandler {
	return &WebsocketHandler{
		channel: make(chan interface{}),
	}
}

func (t *WebsocketHandler) Sink(path string, data interface{}) {
	t.channel <- data
}

func (t *WebsocketHandler) HttpHandler() http.Handler {
	return websocket.Handler(t.ServeWS)
}

func (t WebsocketHandler) ServeWS(ws *websocket.Conn) {
	defer ws.Close()

	path := ws.Request().URL.Path
	log.Printf("Handling [%s]\n", path)

	for payload := range t.channel {
		b, err := json.Marshal(payload)
		if err != nil {
			handleWSError(ws, err)
			return
		}
		if _, err := ws.Write(b); err != nil {
			handleWSError(ws, err)
			return
		}
	}
}

func handleWSError(ws *websocket.Conn, err error) {
	if _, err := ws.Write([]byte(err.Error())); err != nil {
		log.Printf("Error writing error to ws %v", err)
	}
}

/*
	receiver := make(chan interface{})
	connectionId := uuid.New()
	if err := router.AddRoute(path, connectionId, receiver); err != nil {
		err := fmt.Errorf("Invalid path [%s]", path)
		log.Printf(err.Error())
		handleError(ws, err)
		return
	}

	log.Printf("Starting connection to %s...", path)
	defer log.Printf("Stopping connection to %s...", path)
	defer func() {
		router.RemoveRoute(path, connectionId)
		for _ = range receiver {
		}
	}()

	for payload := range receiver {
		b, err := json.Marshal(payload)
		if err != nil {
			handleError(ws, err)
		}
		if _, err := ws.Write(b); err != nil {
			handleError(ws, err)
		}
	}
}
*/
