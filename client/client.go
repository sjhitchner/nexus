package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sjhitchner/infosphere/router"
	"github.com/sjhitchner/infosphere/server"
	"log"
	"net/http"
	"time"
)

func main() {
	for {
		select {
		case <-time.After(5 * time.Second):
			payload := router.Payload{
				Type:      "test",
				Data:      "data",
				CreatedAt: time.Now(),
			}

			b, _ := json.Marshal(payload)

			url := fmt.Sprintf("%s", "http://localhost:8080/v1/put")

			req, err := http.NewRequest(server.METHOD_PUT, url, bytes.NewReader(b))
			if err != nil {
				log.Println(err)
				continue
			}

			req.Header.Add("Content-Type", server.CONTENT_TYPE_JSON)

			client := &http.Client{
			//CheckRedirect: http.redirectPolicyFunc,
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				continue
			}

			if resp.StatusCode != 200 {
				log.Println(err)
				continue
			}
		}
	}
}
