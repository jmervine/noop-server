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

	var store *records.RecordMap
	if record {
		// Use default store.
		store = records.GetStore()
	}

	record := records.NewRecord(r, store)

	f := format
	if record.Echo {
		f = new(formatter.Echo)
	}

	out := recorder.StdRecorder{}
	out.SetFormatter(f)
	out.SetWriter(w)

	w.WriteHeader(record.Status)

	defer func() {
		defer r.Body.Close()

		f := formatter.NewLogFormatter("server.handlerFunc", time.Since(begin), r.Body, verbose)
		log.Printf("%s", f.FormatRecord(record))
	}()

	record.DoSleep() // Only sleeps if sleep is set
	out.WriteOne(record)
}
