package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmervine/env"
)

func init() {
	// set default port if not set
	env.GetOrSetInt("PORT", 3000)
	env.GetOrSetString("ADDR", "0.0.0.0")

	app := os.Getenv("APP_NAME")
	if app == "" {
		app = "noop-server"
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
			if env.GetBool("VERBOSE") {
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

	log.Printf("at=main on=startup addr=%s port=%s\n", env.Get("ADDR"), env.Get("PORT"))
	log.Fatal(http.ListenAndServe(env.Get("ADDR")+":"+env.Get("PORT"), nil))
}
