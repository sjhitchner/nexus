package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleConnection(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(handleConnection))
	defer ts.Close()

	payload := Payload{
		Type:      "test",
		Data:      "data",
		CreatedAt: time.Now(),
	}

	b, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s", ts.URL)

	req, err := http.NewRequest(SUPPORTED_METHOD, url, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", CONTENT_TYPE_JSON)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatal(err)
	}
}
