package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const STATUS_HEADER = "X-HTTP-Status"

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()
	status := http.StatusOK

	defer func() {
		logPrefix := fmt.Sprintf("on=server.handlerFunc method=%s path=%s", r.Method, r.URL.Path)
		log.Printf("%s status=%d took=%v\n", logPrefix, status, time.Since(begin))

		if verbose {
			log.Printf("%s headers:\n%s", logPrefix, r.Header)

			body := &bytes.Buffer{}
			if _, err := io.Copy(body, r.Body); err == nil {
				log.Printf("%s body:\n%s", logPrefix, body.String())
			}
		}
	}()

	if h := r.Header.Get(STATUS_HEADER); h != "" {
		if i, e := strconv.ParseInt(h, 10, 16); e == nil {
			status = int(i)
		} else {
			status = 500
		}
	}

	http.Error(w, fmt.Sprintf("%d %s", status, http.StatusText(status)), status)
}
