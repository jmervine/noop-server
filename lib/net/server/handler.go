package server

import (
	"log"
	"net/http"
	"time"

	"github.com/jmervine/noop-server/lib/recorder"
	"github.com/jmervine/noop-server/lib/records"
	"github.com/jmervine/noop-server/lib/records/formatter"
)

const FLAG_HEADER = "X-NoopServerFlags"

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	record := records.NewRecord(r, record)

	var format formatter.RecordsFormatter
	if record.Echo {
		format = &formatter.Echo{}
	} else {
		format = &formatter.Default{}
	}

	out := recorder.StdRecorder{}
	out.SetFormatter(format)
	out.SetWriter(w)

	w.WriteHeader(record.Status)

	defer func() {
		defer r.Body.Close()

		f := formatter.NewLogFormatter("server.handlerFunc", time.Since(begin), r.Body, verbose)
		log.Printf(f.FormatRecord(record))
	}()

	record.DoSleep() // Only sleeps if sleep is set
	out.WriteOne(record)
}
