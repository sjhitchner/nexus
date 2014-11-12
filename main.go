package main

import (
	"fmt"
	srv "github.com/sjhitchner/library/http"
	agg "github.com/sjhitchner/nexus/interfaces/aggregator"
	handlers "github.com/sjhitchner/nexus/interfaces/handlers"
	"github.com/sjhitchner/nexus/interfaces/multiplex"
	"github.com/sjhitchner/nexus/interfaces/sink"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type LogPublisher struct {
	sync.RWMutex
}

func (t LogPublisher) Publish(b []byte) {
	t.Lock()
	defer t.Unlock()

	fmt.Println("====")
	fmt.Println(string(b))
	fmt.Println("====")
}

var server srv.StoppableServer
var aggregator agg.Aggregator

func main() {
	go signalHandler()

	aggregator = agg.NewAggregator(256, 1, LogPublisher{})
	aggregator.Start()
	//rte = router.NewRouter()
	//rte.AddRoute(agg)
	//rte.AddRoute(agg)
	//rte.AddRoute(agg)
	//srv.AddHandler("application/json", HandleJson)

	multiplexer := multiplex.NewMultiplexer()
	multiplexer.AddSink(sink.LogSink{})
	multiplexer.AddSink(aggregator)

	putHandler := handlers.NewPUTHandler(multiplexer)
	putHandler.AddContentHandler(handlers.CONTENT_TYPE_JSON, handlers.HandleJSONPayload)
	putHandler.AddContentHandler(handlers.CONTENT_TYPE_MSGPACK, handlers.HandleMsgPackPayload)

	apiHandler := handlers.APIHandler{}

	server = srv.NewStoppableServer()
	server.AddHandler("/api", apiHandler)
	server.AddHandler("/v1/put", putHandler)
	server.AddHandler("/ping", handlers.HealthCheckHandler{})
	server.AddHandler("/ws", handlers.WebsocketHandler{})
	if err := server.Start(8080, "static"); err != nil {
		log.Fatal(err)
	}
}

// Handles incoming signals
func signalHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT)

	select {
	case <-ch:
		log.Printf("[INFO] Got a SIGHUP, shutting down service.")
		server.AcceptConnections(false)
		server.Shutdown()
	}
}

/*
func HandleJson(router router.Router, resp http.ResponseWriter, req *http.Request) error {

	body := req.Body
	if body == nil {
		return fmt.Errorf("body is empty")
	}

	dec := json.NewDecoder(body)
	for {
		var payload Payload

		if err := dec.Decode(&payload); err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return err
		}
		router.Route(payload)
	}

	return nil
}
*/
