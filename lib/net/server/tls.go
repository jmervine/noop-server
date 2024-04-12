package server

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"time"
)

func addMTLSSupportToServer(server *http.Server, ca string) {
	// Add the cert chain as the intermediate signs both the servers and the clients certificates
	clientCACert := []byte(ca)

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                clientCertPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	}

	server.TLSConfig = tlsConfig
}

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

func listenAndServeWithTls(server *http.Server, cert string, key string) error {
	crt, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return err
	}

	config := &tls.Config{}
	config.NextProtos = []string{"http/1.1"}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0] = crt

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)},
		config)

	return server.Serve(tlsListener)
}
