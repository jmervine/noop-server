package server

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmervine/noop-server/lib/records"
)

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlerFunc))

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Errorf("Expected nil, got: %v", err)
		return // avoid panic when resp is nil
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got: %d", resp.StatusCode)
	}

	if resp.Request.Method != "GET" {
		t.Errorf("Expected GET, got: %s", resp.Request.Method)
	}
}

func BenchmarkGet(b *testing.B) {
	old := log.Writer()
	log.SetOutput(io.Discard)
	defer func() {
		log.SetOutput(old)
	}()

	server := httptest.NewServer(http.HandlerFunc(handlerFunc))

	for n := 0; n < b.N; n++ {
		_, err := http.Get(server.URL)
		if err != nil {
			b.Errorf("Expected nil, got: %v", err)
		}
	}
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlerFunc))

	resp, err := http.Post(server.URL, "text/html", nil)
	if err != nil {
		t.Errorf("Expected nil, got: %v", err)
		return // avoid panic when resp is nil
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got: %d", resp.StatusCode)
	}

	if resp.Request.Method != "POST" {
		t.Errorf("Expected GET, got: %s", resp.Request.Method)
	}

}

func TestStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlerFunc))
	client := &http.Client{}

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Errorf("Expected nil, got: %v", err)
	}

	req.Header.Add(records.RECORD_HEADER, "status:201")

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Expected nil, got: %v", err)
		return // avoid panic when resp is nil
	}

	if resp.StatusCode != 201 {
		t.Errorf("Expected 201, got: %d", resp.StatusCode)
	}
}
