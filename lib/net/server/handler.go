package server

import (
	"log"
	"net/http"
	"time"

	"github.com/jmervine/noop-server/lib/records"
	"github.com/jmervine/noop-server/lib/records/formatter"
)

const FLAG_HEADER = "X-NoopServerFlags"

func handlerFunc(serverProc int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		begin := time.Now()
		record := records.NewRecord(r, cfg.Listener(), cfg.Sleep, cfg.Echo)

		respFmt := defaultFmt
		if record.Echo {
			respFmt = new(formatter.Echo)
		}

		respBody := respFmt.FormatRecord(record)

		if responders != nil {
			endpoint := r.URL.Path
			if r, ok := responders.Match(endpoint); ok {
				if cfg.Verbose {
					log.Printf("%#v\n", r)
				}

				if r.Status > 0 {
					record.Status = r.Status
				}

				if r.Sleep > 0 {
					record.Sleep = time.Duration(r.Sleep * uint(time.Millisecond))
				}

				if r.Response != "" {
					respBody = r.Response
				}
			}

		}

		if store != nil {
			store.Add(record)
		}

		// Stream record to record file, if stream enabled
		if cfg.StreamRecord {
			stream.WriteOne(record)
			stream.WriteString("\n")
		}

		record.DoSleep() // Only sleeps if sleep is set
		http.Error(w, respBody, record.Status)

		logFmt := formatter.NewLogFormatter("server.handlerFunc", serverProc, time.Since(begin), r.Body, cfg.Verbose)
		log.Printf("%s", logFmt.FormatRecord(record))
	}
}
