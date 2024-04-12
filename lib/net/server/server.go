package server

import (
	"net/http"

	"github.com/jmervine/noop-server/lib/config"
)

var verbose bool

func Start(c *config.Config) error {
	verbose = c.Verbose

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunc)

	svr := &http.Server{Addr: c.Listener(), Handler: mux}
	if c.MTLSEnabled() {
		addMTLSSupportToServer(svr, c.CertCAPath)
	}

	if c.TLSEnabled() {
		return listenAndServeWithTls(svr, c.CertPrivatePath, c.CertKeyPath)
	}

	return svr.ListenAndServe()
}
