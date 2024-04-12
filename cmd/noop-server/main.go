package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmervine/noop-server/lib/config"
	"github.com/jmervine/noop-server/lib/net/handler"
)

var (
	cfg  *config.Config
	cert tls.Certificate
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

func main() {
	cfg = config.Init(os.Args)

	log.SetPrefix(fmt.Sprintf("app=%s ", cfg.App))
	log.Printf("on=startup %s\n", cfg.ToString())

	var err error
	if cfg.TLSEnabled() {
		if cert, err = tls.X509KeyPair([]byte(cfg.CertPrivatePath), []byte(cfg.CertKeyPath)); err != nil {
			log.Fatalf("error=\"%v\" cert=\"%s\" key=\"%s\"\n", err, cfg.CertPrivatePath, cfg.CertKeyPath)
		}
	}

	mux := http.NewServeMux()

	handlerFunc := handler.Init(cfg)
	mux.HandleFunc("/", handlerFunc)

	server := &http.Server{Addr: cfg.Listener(), Handler: mux}
	if cfg.MTLSEnabled() {
		addMTLSSupportToServer(server)
	}

	if cfg.TLSEnabled() {
		log.Fatal(listenAndServeTLSKeyPair(server, cert))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}
