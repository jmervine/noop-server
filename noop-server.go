package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmervine/env"
	log "github.com/jmervine/readable"
)

func init() {
	// set default port if not set
	env.GetOrSetInt("PORT", 3000)

	log.SetPrefix("[noop-server]:")
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	log.SetDebug(env.GetBool("VERBOSE"))
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		status := http.StatusOK

		defer func() {
			log.Print("method", r.Method, "path", r.URL.Path, "status", status, "took", time.Since(begin))
			if env.GetBool("VERBOSE") {
				log.Debug("method", r.Method, "path", r.URL.Path, "headers", r.Header)
				body, err := ioutil.ReadAll(r.Body)
				if err == nil {
					log.Debug("method", r.Method, "path", r.URL.Path, "body", string(body))
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

	log.Print("at=main on=startup addr", env.Get("ADDR"), "port", env.Get("PORT"))
	log.Fatal(http.ListenAndServe(env.Get("ADDR")+":"+env.Get("PORT"), nil))
}
