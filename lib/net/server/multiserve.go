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

// Support passing server in for testing, default is to pass nil
func multiListenAndServe(server *http.Server, s time.Duration) error {
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

	if cfg.MaxProcs() == 1 {
		log.Println(server.Serve(listener))
		return nil
	}

	var wg sync.WaitGroup
	for i := 0; i < cfg.MaxProcs(); i++ {
		wg.Add(1)
		go func(i int) {
			if server == nil {
				server = buildServer(i, s)
			}

			log.Printf("at=server.Start in=server.multiListenAndServe listener=%03d\n", i)
			log.Println(server.Serve(listener))
			wg.Done()
		}(i)
	}
	wg.Wait()

	return nil
}

// Server has access to configurations via 'cfg'
func buildServer(n int, s time.Duration) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunc(n))

	server := &http.Server{
		Addr:    cfg.Listener(),
		Handler: mux,

		// Short timeouts -- this should always be fast.
		// Really these timeouts should be closer to 100us
		ReadTimeout:       (25 * time.Millisecond) + s,
		WriteTimeout:      (25 * time.Millisecond) + s,
		IdleTimeout:       (25 * time.Millisecond) + s,
		ReadHeaderTimeout: (25 * time.Millisecond) + s,
	}

	if cfg.MTLSEnabled() {
		addMTLSSupportToServer(server, cfg.CertCAPath)
	}

	return server
}
