package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jmervine/noop-server/lib/recorder"
)

const FLAG_HEADER = "X-NoopServerFlags"

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	record := recorder.NewRecord(r, record)

	defer func() {
		logPrefix := fmt.Sprintf("on=server.handlerFunc method=%s path=%s", r.Method, r.URL.Path)
		log.Printf("%s status=%d took=%v\n", logPrefix, record.Status, time.Since(begin))

		if verbose {
			log.Printf("%s headers:\n%s", logPrefix, r.Header)

			body := &bytes.Buffer{}
			if _, err := io.Copy(body, r.Body); err == nil {
				log.Printf("%s body:\n%s", logPrefix, body.String())
			}
		}
	}()

	record.DoSleep() // Only sleeps if sleep is set

	body := record.EchoString()
	http.Error(w, body, record.Status)
}
