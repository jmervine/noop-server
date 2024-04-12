package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmervine/noop-server/lib/config"
)

var (
	cfg  config.Config
	cert tls.Certificate
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

func handler(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()
	status := http.StatusOK

	defer func() {
		logPrefix := fmt.Sprintf("on=http.HandleFunc method=%s path=%s", r.Method, r.URL.Path)

		log.Printf("%s status=%d took=%v\n", logPrefix, status, time.Since(begin))

		if cfg.Verbose {
			log.Printf("%s headers:\n%s", logPrefix, r.Header)

			body := &bytes.Buffer{}
			if _, err := io.Copy(body, r.Body); err == nil {
				log.Printf("%s body:\n%s", logPrefix, body.String())
			}
		}
	}()

	if h := r.Header.Get("X-HTTP-Status"); h != "" {
		if i, e := strconv.ParseInt(h, 10, 16); e == nil {
			status = int(i)
		} else {
			status = 500
		}
	}

	http.Error(w, fmt.Sprintf("%d %s", status, http.StatusText(status)), status)
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
	mux.HandleFunc("/", handler)

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
