package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type APIHandler struct {
}

func (t APIHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Println("APIHandler")

	/*
		/filter
		/logging
		/log
		/status

	*/

	resp.WriteHeader(200)
}

type HealthCheckHandler struct {
}

func (t HealthCheckHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Println("PING!")
	resp.WriteHeader(200)
}

func handleError(resp http.ResponseWriter, statusCode int, msg string, args ...interface{}) {
	resp.WriteHeader(statusCode)
	log.Println("ERROR:", fmt.Sprintf(msg, args...))
	fmt.Fprintf(resp, msg, args)
}

/*
func (t *server) handleConnection(resp http.ResponseWriter, req *http.Request) {
	log.Println("request")



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
*/
