package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmervine/noop-server/Godeps/_workspace/src/github.com/jmervine/env"
	log "github.com/jmervine/noop-server/Godeps/_workspace/src/github.com/jmervine/readable"
	"github.com/jmervine/noop-server/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
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
	router := httprouter.New()

	methods := []string{
		"DELETE",
		"GET",
		"HEAD",
		"POST",
		"OPTIONS",
		"PUT",
	}

	for _, m := range methods {
		router.Handle(m, "/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

			http.Error(w, http.StatusText(200), 200)
		})

		router.Handle(m, "/:status", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			begin := time.Now()
			status := http.StatusOK

			defer func(s int) {
				log.Print("method", r.Method, "path", r.URL.Path, "status", s, "took", time.Since(begin))
				if env.GetBool("VERBOSE") {
					log.Debug("method", r.Method, "path", r.URL.Path, "headers", r.Header)
					body, err := ioutil.ReadAll(r.Body)
					if err == nil {
						log.Debug("method", r.Method, "path", r.URL.Path, "body", string(body))
					}
				}
			}(status)

			i, _ := strconv.ParseInt(p.ByName("status"), 10, 16)
			if i > 99 {
				status = int(i)
			}

			http.Error(w, http.StatusText(status), status)
		})
	}

	log.Print("at=main on=startup addr", env.Get("ADDR"), "port", env.Get("PORT"))
	log.Fatal(http.ListenAndServe(env.Get("ADDR")+":"+env.Get("PORT"), router))

}
