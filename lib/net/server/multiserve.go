package server

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// From: https://gist.github.com/shivakar/cd52b5594d4912fbeb46
// ---------------------------------------------------------------------------
// From https://golang.org/src/net/http/server.go
// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func multiListenAndServe() error {
	ln, err := net.Listen("tcp", cfg.Listener())
	if err != nil {
		return err
	}

	serverConfig, err := tlsConfig()
	if err != nil {
		return err
	}

	var listener net.Listener
	listener = tcpKeepAliveListener{ln.(*net.TCPListener)}
	if cfg.TLSEnabled() {
		listener = tls.NewListener(listener, serverConfig)
	}

	var wg sync.WaitGroup
	for i := 0; i < cfg.MaxProcs(); i++ {
		wg.Add(1)
		go func(i int) {
			server := buildServer(i)

			log.Printf("at=server.Start in=server.multiListenAndServe listener=%03d\n", i)
			log.Println(server.Serve(listener))
			wg.Done()
		}(i)
	}
	wg.Wait()

	return nil
}

// Server has access to configurations via 'cfg'
func buildServer(n int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunc(n))

	server := &http.Server{
		Addr:    cfg.Listener(),
		Handler: mux,

		// Short timeouts -- this should always be fast.
		// Really these timeouts should be closer to 100us
		ReadTimeout:       10 * time.Millisecond,
		WriteTimeout:      10 * time.Millisecond,
		IdleTimeout:       10 * time.Millisecond,
		ReadHeaderTimeout: 10 * time.Millisecond,
	}

	if cfg.MTLSEnabled() {
		addMTLSSupportToServer(server, cfg.CertCAPath)
	}

	return server
}
