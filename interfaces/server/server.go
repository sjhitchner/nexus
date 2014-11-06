package server

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	. "github.com/sjhitchner/infosphere/domain"
	"github.com/sjhitchner/infosphere/interfaces/multiplex"
	"log"
	"net/http"
	_ "net/http/pprof"
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

type ContentHandler func(multiplexer multiplexer.Multiplexer, resp http.ResponseWriter, req *http.Request) error

type Server interface {
	AddHandler(contentType string, handler ContentHandler)
	AcceptConnections(accept bool)
	Start(port int, staticPath string) error
	Shutdown()
}

type server struct {
	sync.RWMutex
	port              int
	router            router.Router
	handlers          ContentHandlerMap
	listener          *StoppableListener
	acceptConnections bool
}

func NewServer(multiplexer multiplexer.Multiplexer) Server {
	return &server{
		multiplexer:       multiplexer,
		handlers:          make(ContentHandlerMap),
		listener:          nil,
		acceptConnections: true,
	}
}

func (t *server) AddHandler(contentType string, handler ContentHandler) {
	t.Lock()
	defer t.Unlock()
	t.handlers[contentType] = handler
}

func (t *server) Start(port int, staticPath string) error {
	http.HandleFunc("/v1/put", t.handleConnection)
	http.HandleFunc("/ping", t.handlePing)
	http.Handle("/", http.FileServer(http.Dir(staticPath)))
	http.Handle("/websocket/", websocket.Handler(websocketHandler))

	var err error
	t.listener, err = NewListener(fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return Serve(t.listener)
}

func (t *server) AcceptConnections(accept bool) {
	log.Println("Not accepting connections")
	t.acceptConnections = accept
}

func (t *server) Shutdown() {
	t.acceptConnections = false
	t.listener.Stop()
}

func (t *server) handleConnection(resp http.ResponseWriter, req *http.Request) {
	log.Println("request")

	if req.Method != METHOD_PUT {
		handleError(resp, http.StatusMethodNotAllowed, "Invalid HTTP Verb [%s]", req.Method)
		return
	}

	contentType := strings.ToLower(req.Header.Get("Content-Type"))
	t.RLock()
	handlerFunc, ok := t.handlers[contentType]
	t.RUnlock()
	if !ok {
		handleError(resp, http.StatusUnsupportedMediaType, "Invalid Content-Type=[%s]", contentType)
		return
	}

	if err := handlerFunc(t.multiplexer, resp, req); err != nil {
		handleError(resp, http.StatusInternalServerError, "error occurred")
		return
	}

	resp.WriteHeader(http.StatusOK)
	return
}

// Health Check
func (t *server) handlePing(resp http.ResponseWriter, req *http.Request) {
	if req.Method != METHOD_GET {
		handleError(resp, http.StatusMethodNotAllowed, "Invalid HTTP Verb [%s]", req.Method)
		return
	}

	statusCode := http.StatusOK
	if !t.acceptConnections {
		statusCode = http.StatusInternalServerError
		log.Println("Out of service")
	}

	resp.WriteHeader(statusCode)
	return
}

func handleError(resp http.ResponseWriter, statusCode int, msg string, args ...interface{}) {
	resp.WriteHeader(statusCode)
	fmt.Fprintf(resp, msg, args)
}
