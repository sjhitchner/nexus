package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sjhitchner/nexus/interfaces/multiplex"
	"io"
	"net/http"
	"strings"
	"sync"
)

const (
	METHOD_PUT           = "PUT"
	METHOD_GET           = "GET"
	CONTENT_TYPE_JSON    = "application/json"
	CONTENT_TYPE_MSGPACK = "application/msgpack"
)

type ContentHandlerMap map[string]ContentHandler

type ContentHandler func(multiplexer multiplex.Multiplexer, resp http.ResponseWriter, req *http.Request) error

type PUTHandler struct {
	sync.RWMutex
	handlers    ContentHandlerMap
	multiplexer multiplex.Multiplexer
}

func NewPUTHandler(multiplexer multiplex.Multiplexer) *PUTHandler {
	return &PUTHandler{
		handlers:    make(ContentHandlerMap),
		multiplexer: multiplexer,
	}
}

func (t *PUTHandler) AddContentHandler(contentType string, handler ContentHandler) {
	t.Lock()
	defer t.Unlock()
	t.handlers[contentType] = handler
}

func (t PUTHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != METHOD_PUT {
		http.Error(resp, fmt.Sprintf("Invalid HTTP Verb [%s]", req.Method), http.StatusMethodNotAllowed)
		return
	}

	contentType := strings.ToLower(req.Header.Get("Content-Type"))
	t.RLock()
	handlerFunc, ok := t.handlers[contentType]
	t.RUnlock()
	if !ok {
		http.Error(resp, fmt.Sprintf("Invalid Content-Type=[%s]", contentType), http.StatusUnsupportedMediaType)
		return
	}

	if err := handlerFunc(t.multiplexer, resp, req); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	return
}

func HandleJSONPayload(multiplexer multiplex.Multiplexer, resp http.ResponseWriter, req *http.Request) error {
	// TODO: whitelist / blacklist paths?
	path, err := parsePathFromRequestURI(req.RequestURI)
	if err != nil {
		return err
	}

	body := req.Body
	if body == nil {
		return fmt.Errorf("body is empty")
	}

	dec := json.NewDecoder(body)
	for {
		var data interface{}

		if err := dec.Decode(&data); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		multiplexer.Multiplex(path, data)
	}

	return nil
}

func HandleMsgPackPayload(multiplexer multiplex.Multiplexer, resp http.ResponseWriter, req *http.Request) error {
	path, err := parsePathFromRequestURI(req.RequestURI)
	if err != nil {
		return err
	}

	body := req.Body
	if body == nil {
		return fmt.Errorf("body is empty")
	}

	dec := json.NewDecoder(body)
	for {
		var data interface{}

		if err := dec.Decode(&data); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		multiplexer.Multiplex(path, data)
	}

	return nil
}

func parsePathFromRequestURI(uri string) (string, error) {
	arr := strings.SplitAfterN(uri, "/", 4)
	if len(arr) != 4 {
		return "", fmt.Errorf("Invalid requestURI %s", uri)
	}
	return arr[3], nil
}
