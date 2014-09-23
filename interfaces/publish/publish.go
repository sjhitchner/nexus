package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var shutdownChannels []chan bool
var queue = make(chan Payload, 10000)
var wg = sync.WaitGroup{}

func Publish(data Payload) {
	queue <- data
}

func StartPublisher(threads int) { //, publishers ...Publisher) error {
	go signalHandler()

	shutdownChannel = make([]chan bool, threads)

	for i := 0; i < threads; i++ {

		shutdownChannels[i] = make(chan bool)
		go publishWorker(shutdownChannels[i])
	}

	wg.Wait()
}

func ShutdownPublisher() {
	for i := 0; i < len(shutdownChannels); i++ {
		log.Println("X")
		shutdownChannels[i] <- true
	}

}

func publishWorker(wg sync.WaitGroup, shutdownChannel chan bool) {
	log.Println("starting...")
	wg.Add(1)

	for {
		select {
		case payload := <-publishQueue:
			log.Println(payload)
		case <-shutdownChannel:
			log.Println("exiting...")
			wg.Done()
			return
		}
	}
}

/*
type Publisher interface {
	Publish(data Payload)
}

type BatchPublisher struct {
}

func (t *BatchPublisher) Publish(data Payload) {

}
*/

// Handles incoming signals
func signalHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT)

	// Quit the workers if we get signalled
	select {
	case <-ch:
		log.Printf("[INFO] Got a SIGHUP, quiting all workers.")
		ShutdownPublisher()
	}
}
