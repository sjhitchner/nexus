package main

import (
	//"github.com/crowdmob/goamz/aws"
	srv "github.com/sjhitchner/library/http"
	agg "github.com/sjhitchner/nexus/interfaces/aggregator"
	handlers "github.com/sjhitchner/nexus/interfaces/handlers"
	"github.com/sjhitchner/nexus/interfaces/multiplex"
	"github.com/sjhitchner/nexus/interfaces/publish"
	//"github.com/sjhitchner/nexus/interfaces/sink"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	//AWS_CREDENTIALS = aws.Auth{"zzzz", "zzzz"}
	BUCKET           = "initium-logs"
	ENV_NEXUS_STATUS = "NEXUS_STATIC"
)

var server srv.StoppableServer
var aggregator agg.Aggregator

func main() {
	go signalHandler()

	//s3publisher := publish.NewS3Publisher(AWS_CREDENTIALS, aws.APNortheast, BUCKET)

	aggregator = agg.NewAggregator(512, time.Minute/4, publish.LogPublisher{})
	//aggregator = agg.NewAggregator(256, 2, s3publisher)
	//aggregator.Start()
	//rte = router.NewRouter()
	//rte.AddRoute(agg)
	//rte.AddRoute(agg)
	//rte.AddRoute(agg)
	//srv.AddHandler("application/json", HandleJson)
	wsHandler := handlers.NewWebsocketHandler()

	multiplexer := multiplex.NewMultiplexer()
	//multiplexer.AddSink(sink.LogSink{})
	multiplexer.AddSink(aggregator)
	multiplexer.AddSink(wsHandler)

	putHandler := handlers.NewPUTHandler(multiplexer)
	putHandler.AddContentHandler(handlers.CONTENT_TYPE_JSON, handlers.HandleJSONPayload)
	putHandler.AddContentHandler(handlers.CONTENT_TYPE_MSGPACK, handlers.HandleMsgPackPayload)

	apiHandler := handlers.APIHandler{}

	server = srv.NewStoppableServer()
	server.AddHandler("/api", apiHandler)
	server.AddHandler("/v1/put/", putHandler)
	server.AddHandler("/ping", handlers.HealthCheckHandler{})
	server.AddHandler("/ws", wsHandler.HttpHandler())
	server.AddHandler("/", http.FileServer(http.Dir(GetStaticDirectory())))
	if err := server.Start(8080); err != nil {
		log.Fatal(err)
	}
}

func GetStaticDirectory() string {
	staticDir := os.Getenv(ENV_NEXUS_STATUS)
	if staticDir == "" {
		panic(fmt.Errorf("Need to set [%s]", ENV_NEXUS_STATUS))
	}
	fileInfo, err := os.Stat(staticDir)
	if err != nil {
		panic(fmt.Errorf("Error with [%s] %v", ENV_NEXUS_STATUS, err))
	}
	if !fileInfo.IsDir() {
		panic(fmt.Errorf("Error with [%s] not valid directory", ENV_NEXUS_STATUS))
	}
	return staticDir
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

		// TODO: this is bad shouldn't use a time here
		//time.Sleep(1 * time.Second)
		aggregator.Shutdown()
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
