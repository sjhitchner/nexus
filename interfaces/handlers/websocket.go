package handlers

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
)

type WebsocketHandler struct {
}

func NewWebsocketHandler() *WebsocketHandler {
	return &WebsocketHandler{}
}

func (t *WebsocketHandler) HttpHandler() http.Handler {
	return websocket.Handler(t.ServeWS)
}

func (t WebsocketHandler) ServeWS(ws *websocket.Conn) {
	path := ws.Request().URL.Path
	log.Printf("Handling [%s]\n", path)
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
