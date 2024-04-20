package server

import (
	"log"
	"net/http"
	"time"

	"github.com/jmervine/noop-server/lib/records"
	"github.com/jmervine/noop-server/lib/records/formatter"
)

const FLAG_HEADER = "X-NoopServerFlags"

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	begin := time.Now()
	record := records.NewRecord(r, cfg.Listener())
	if store != nil {
		store.Add(record)
	}

	respFmt := defaultFmt
	if record.Echo {
		respFmt = new(formatter.Echo)
	}

	// Stream record to record file, if stream enabled
	if cfg.StreamRecord {
		stream.WriteOne(record)
		stream.WriteString("\n")
	}

	record.DoSleep() // Only sleeps if sleep is set
	http.Error(w, respFmt.FormatRecord(record), record.Status)

	logFmt := formatter.NewLogFormatter("server.handlerFunc", time.Since(begin), r.Body, cfg.Verbose)
	log.Printf("%s", logFmt.FormatRecord(record))
}
