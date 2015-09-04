package main

import (
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
				log.Print("method", m, "path", r.URL.Path, "status", status, "took", time.Since(begin))
			}()

			http.Error(w, http.StatusText(200), 200)
		})

		router.Handle(m, "/:status", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			begin := time.Now()
			status := http.StatusOK

			defer func(s int) {
				log.Print("method", m, "path", r.URL.Path, "status", s, "took", time.Since(begin))
			}(status)

			i, _ := strconv.ParseInt(p.ByName("status"), 10, 16)
			if i > 99 {
				status = int(i)
			}

			http.Error(w, http.StatusText(status), status)
		})
	}

	log.Print("at=main on=startup addr", env.Get("ADDR"), "port", env.Get("PORT"))
	log.Fatal(http.ListenAndServe(env.GetString("ADDR")+":"+env.GetString("PORT"), router))

}
