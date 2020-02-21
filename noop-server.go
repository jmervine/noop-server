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
	port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	addr = os.Getenv("ADDR")
	if addr == "" {
		addr = "0.0.0.0"
	}

	app = os.Getenv("APP_NAME")
	if app == "" {
		app = "noop-server"
	}

	if os.Getenv("VERBOSE") != "" {
		verbose = true
	}

	log.SetPrefix("app=" + app + " ")
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		status := http.StatusOK

		defer func() {
			logPrefix := fmt.Sprintf("at=main on=http.HandleFunc method=%s path=%s", r.Method, r.URL.Path)
			log.Printf("%s status=%d took=%v\n", logPrefix, status, time.Since(begin))
			if verbose {
				log.Printf("%s headers:\n%s", r.Header)
				body, err := ioutil.ReadAll(r.Body)
				if err == nil {
					log.Printf("%s body:\n%s", string(body))
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
	})

	log.Printf("at=main on=startup addr=%s port=%s\n", addr, port)
	log.Fatal(http.ListenAndServe(addr+":"+port, nil))
}
