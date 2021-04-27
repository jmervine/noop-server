package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	port, addr, app, listener string
	mtlsCert, tlsCert, tlsKey string
	verbose, mtls, stls       bool
	cert                      tls.Certificate
)

func init() {
	var ok bool

	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	if app, ok = os.LookupEnv("APP_NAME"); !ok {
		app = "noop-server"
	}

	log.SetPrefix(fmt.Sprintf("app=%s ", app))

	if port, ok = os.LookupEnv("PORT"); !ok {
		port = "3000"
	}

	if addr, ok = os.LookupEnv("ADDR"); !ok {
		addr = "0.0.0.0"
	}

	listener = fmt.Sprintf("%s:%s", addr, port)

	if os.Getenv("VERBOSE") != "" {
		verbose = true
	}

	tlsCert = os.Getenv("TLS_CERT")
	tlsKey = os.Getenv("TLS_KEY")
	if tlsCert != "" && tlsKey != "" {
		stls = true

		var e error
		if cert, e = tls.X509KeyPair([]byte(tlsCert), []byte(tlsKey)); e != nil {
			log.Fatalln(e)
		}

		// no mtls with tls
		mtlsCert = os.Getenv("MTLS_CA_CHAIN_CERT")
		if mtlsCert != "" {
			mtls = true
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()
	status := http.StatusOK

	defer func() {
		logPrefix := fmt.Sprintf("on=http.HandleFunc method=%s path=%s", r.Method, r.URL.Path)

		log.Printf("%s status=%d took=%v\n", logPrefix, status, time.Since(begin))
		if verbose {
			log.Printf("%s headers:\n%s", logPrefix, r.Header)
			body, err := ioutil.ReadAll(r.Body)
			if err == nil {
				log.Printf("%s body:\n%s", logPrefix, string(body))
			}
		}
	}()

	statusHeader := r.Header.Get("X-HTTP-Status")
	if statusHeader != "" {
		i, e := strconv.ParseInt(statusHeader, 10, 16)
		if e == nil {
			status = int(i)
		} else {
			status = 500
		}
	}

	http.Error(w, fmt.Sprintf("%d %s", status, http.StatusText(status)), status)
}

func main() {
	log.Printf("on=startup addr=%s port=%s mtls=%v ssl=%v\n", addr, port, mtls, stls)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	server := &http.Server{Addr: listener, Handler: mux}
	if mtls {
		addMTLSSupportToServer(server)
	}

	if !stls {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Fatal(listenAndServeTLSKeyPair(server, cert))
	}
}
