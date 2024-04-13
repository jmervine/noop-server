package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const FLAG_HEADER = "X-NoopServerFlags"

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	// unset default should work, or be handled
	flags := parseHeaderFlags(r.Header.Get(FLAG_HEADER))

	defer func() {
		logPrefix := fmt.Sprintf("on=server.handlerFunc method=%s path=%s", r.Method, r.URL.Path)
		log.Printf("%s status=%d took=%v\n", logPrefix, flags.Status(), time.Since(begin))

		if verbose {
			log.Printf("%s headers:\n%s", logPrefix, r.Header)

			body := &bytes.Buffer{}
			if _, err := io.Copy(body, r.Body); err == nil {
				log.Printf("%s body:\n%s", logPrefix, body.String())
			}
		}
	}()

	flags.Sleep() // Only sleeps if sleep is set
	http.Error(w, flags.Echo(r), flags.Status())
}
