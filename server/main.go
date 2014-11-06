package main

import (
	"encoding/json"
	"fmt"
	. "github.com/sjhitchner/infosphere/domain"
	"github.com/sjhitchner/infosphere/interfaces/aggregator"
	"github.com/sjhitchner/infosphere/interfaces/router"
	"github.com/sjhitchner/infosphere/interfaces/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

/*
Consume Daemon To Aggregate Content

UDP socket, get content channel it, aggregate, publish

Filter, channel, filter datasource

type Filter interface {
}
*/

var plock = sync.RWMutex{}

type LogPublisher struct {
}

func (t LogPublisher) Publish(b []byte) {
	plock.Lock()
	defer plock.Unlock()

	fmt.Println("====")
	fmt.Println(string(b))
	fmt.Println("====")
}

var agg aggregator.Aggregator
var rte router.Router
var srv server.Server

func main() {
	go signalHandler()

	agg = aggregator.NewAggregator(LogPublisher{})
	agg.Start()

	rte = router.NewRouter()
	rte.AddRoute(agg)
	rte.AddRoute(agg)
	rte.AddRoute(agg)

	srv = server.NewServer(rte)
	srv.AddHandler("application/json", HandleJson)
	if err := srv.Start(8080, "static"); err != nil {
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
		srv.AcceptConnections(false)
		agg.Shutdown()
		srv.Shutdown()
	}
}

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
