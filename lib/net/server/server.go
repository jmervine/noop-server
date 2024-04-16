package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmervine/noop-server/lib/config"
	"github.com/jmervine/noop-server/lib/recorder"
	"github.com/jmervine/noop-server/lib/records"
	"github.com/jmervine/noop-server/lib/records/formatter"
)

var verbose bool
var record bool
var format formatter.RecordsFormatter = new(formatter.Default)

func Start(c *config.Config) error {
	verbose = c.Verbose
	record = c.Record

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunc)

	svr := &http.Server{Addr: c.Listener(), Handler: mux}
	if c.MTLSEnabled() {
		addMTLSSupportToServer(svr, c.CertCAPath)
	}

	// TODO: Consider remove all recorder handling in net/server to recorder.
	if record {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)
		go func() {
			// As currently designed SIGHUP will record only, without stopping the
			// server. SIGINT will record and exit.
			sig := <-sigChan

			r := recorder.SerialRecorder{}

			file, err := c.RecordFile()
			if err != nil {
				log.Fatalf("Error creating '%s' with %v\n", c.RecordTarget, err)
			}
			defer file.Close()

			r.SetFormatter(formatter.NoopClient{})
			r.SetWriter(file)

			store := records.GetStore()
			log.Printf("on=server.Start record-target='%s' record-count=%d\n", c.RecordTarget, store.Size())
			r.WriteAll(store)

			if sig == syscall.SIGINT {
				os.Exit(0)
			}
		}()
	}

	if c.TLSEnabled() {
		return listenAndServeWithTls(svr, c.CertPrivatePath, c.CertKeyPath)
	}

	return svr.ListenAndServe()
}
