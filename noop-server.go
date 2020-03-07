package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	port, addr, app string
	verbose         bool
)

func init() {
    var ok bool
    if port, ok = os.LookupEnv("PORT"); !ok {
		port = "3000"
	}

	if addr, ok = os.LookupEnv("ADDR"); !ok {
		addr = "0.0.0.0"
	}

	if app, ok = os.LookupEnv("APP_NAME"); !ok {
		app = "noop-server"
	}

	if os.Getenv("VERBOSE") != "" {
		verbose = true
	}

	log.SetPrefix("app=" + app + " ")
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

func handler(w http.ResponseWriter, r *http.Request) {
    begin := time.Now()
    status := http.StatusOK

    defer func() {
        logPrefix := fmt.Sprintf("on=http.HandleFunc method=%s path=%s ", r.Method, r.URL.Path)

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
    log.SetFlags(0)
    log.SetPrefix("at=main ")

	http.HandleFunc("/", handler)

	log.Printf("on=startup addr=%s port=%s\n", addr, port)
	log.Fatal(http.ListenAndServe(addr+":"+port, nil))
}
