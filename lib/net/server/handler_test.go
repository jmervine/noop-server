package server

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jmervine/noop-server/lib/config"
	"github.com/jmervine/noop-server/lib/records"
)

var tclient *http.Client

func init() {
	cfg = new(config.Config)
	tclient = &http.Client{
		Timeout: 10 * time.Millisecond,
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlerFunc(0)))
	defer server.Close()

	resp, err := tclient.Get(server.URL)
	if err != nil {
		t.Errorf("Expected nil, got: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got: %d", resp.StatusCode)
	}

	if resp.Request.Method != "GET" {
		t.Errorf("Expected GET, got: %s", resp.Request.Method)
	}
}

func BenchmarkServer_handlerWithGet(b *testing.B) {
	old := log.Writer()
	log.SetOutput(io.Discard)
	defer func() {
		log.SetOutput(old)
	}()

	server := httptest.NewServer(http.HandlerFunc(handlerFunc(0)))

	for n := 0; n < b.N; n++ {
		_, err := http.Get(server.URL)
		if err != nil {
			b.Errorf("Expected nil, got: %v", err)
		}
	}
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlerFunc(0)))
	defer server.Close()

	resp, err := tclient.Post(server.URL, "text/html", nil)
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
	server := httptest.NewServer(http.HandlerFunc(handlerFunc(0)))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Errorf("Expected nil, got: %v", err)
	}

	req.Header.Add(records.RECORD_HEADER, "status:201")

	resp, err := tclient.Do(req)
	if err != nil {
		t.Errorf("Expected nil, got: %v", err)
		return // avoid panic when resp is nil
	}

	if resp.StatusCode != 201 {
		t.Errorf("Expected 201, got: %d", resp.StatusCode)
	}
}
