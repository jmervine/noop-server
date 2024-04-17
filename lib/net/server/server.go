package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmervine/noop-server/lib/config"
	"github.com/jmervine/noop-server/lib/recorder"
	"github.com/jmervine/noop-server/lib/records"
	"github.com/jmervine/noop-server/lib/records/formatter"
	"github.com/pkg/errors"
)

var cfg *config.Config

// Setup (mostly) empty defaults
var defaultFmt formatter.RecordsFormatter = formatter.Default{}
var store *records.RecordMap

// For record writer, either stream-record or record
var stream recorder.Recorder

func Start(c *config.Config) error {
	cfg = c

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunc)

	svr := &http.Server{
		Addr:    c.Listener(),
		Handler: mux,

		// Short timeouts -- this should always be fast.
		// Really these timeouts should be closer to 100us
		ReadTimeout:       10 * time.Millisecond,
		WriteTimeout:      10 * time.Millisecond,
		IdleTimeout:       10 * time.Millisecond,
		ReadHeaderTimeout: 10 * time.Millisecond,
	}

	if c.MTLSEnabled() {
		addMTLSSupportToServer(svr, c.CertCAPath)
	}

	// TODO: Consider remove all recorder handling in net/server to recorder.
	if c.Recording() {
		store = records.GetStore()

		file, err := c.RecordFile()
		if err != nil {
			return errors.Errorf("Error creating '%s' with %v\n", c.RecordTarget, err)
		}
		defer file.Close()

		stream = &recorder.StdRecorder{}
		format := formatter.NoopClient{}

		// TODO: Make StdFormatter#WriteHeader function to handle this for all things.
		if c.StreamRecord {
			format.NewLine = true

			head := fmt.Sprintf("# Stream started: %s\n", time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
			_, err := file.WriteString(head)
			if err != nil {
				return errors.Errorf("error writing to stream file: %v", err)
			}

			err = file.Sync()
			if err != nil {
				return errors.Errorf("error syncing to stream file: %v", err)
			}

			// To ensure things are flushed correctly
			defer file.Sync()
		}

		// Call after writing the header to make sure 'file' is available
		stream.SetFormatter(format)
		stream.SetWriter(file)
	}

	if c.Record {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)
		go func() {
			// As currently designed SIGHUP will record only, without stopping the
			// server. SIGINT will record and exit.
			sig := <-sigChan

			store := records.GetStore()
			log.Printf("on=server.Start record-target='%s' record-count=%d\n", c.RecordTarget, store.Size())
			stream.WriteAll(store)

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
