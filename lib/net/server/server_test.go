package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/jmervine/noop-server/lib/config"
)

func setupBenchmark(b *testing.B, server *http.Server) func() {
	cfg = &config.Config{
		App:    "test-noop-server",
		Port:   "3333",
		Addr:   "localhost",
		NProcs: 4,
	}

	old := log.Writer()
	log.SetOutput(io.Discard)

	go func() {
		err := multiListenAndServe(server, 0)
		if err != nil {
			b.Errorf("Expected nil, got: %v", err)
		}
	}()

	// Give server time to start.
	time.Sleep(5 * time.Second)

	return func() {
		server.Shutdown(context.Background())
		log.SetOutput(old)
		cfg = nil
	}
}

func BenchmarkServer(b *testing.B) {
	server := buildServer(0, 0)
	closer := setupBenchmark(b, server)
	defer closer()

	for n := 0; n < b.N; n++ {
		_, err := http.Get("http://" + cfg.Listener())
		if err != nil {
			b.Errorf("Expected nil, got: %v", err)
		}
	}
}
