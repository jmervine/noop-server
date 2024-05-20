package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmervine/noop-server/lib/config"
	"github.com/jmervine/noop-server/lib/recorder"
	"github.com/jmervine/noop-server/lib/records"
	"github.com/jmervine/noop-server/lib/records/formatter"
	"github.com/jmervine/noop-server/lib/responder"
	"github.com/pkg/errors"
)

var cfg *config.Config

// Setup (mostly) empty defaults
var defaultFmt formatter.RecordsFormatter = formatter.Default{}
var store *records.RecordMap

// For record writer, either stream-record or record
var stream recorder.Recorder
var responders *responder.Responders

func Start(c *config.Config) error {
	cfg = c

	var err error
	if c.Script != "" {
		if responders, err = responder.Load(c.Script); err != nil {
			return err
		}
	}

	if c.Recording() {
		store = records.GetStore()

		file, err := c.RecordFile()
		if err != nil {
			return errors.Errorf("Error creating '%s' with %v\n", c.RecordTarget, err)
		}
		defer file.Close()

		format := cfg.RecordFormatter()
		stream = &recorder.StdRecorder{}
		stream.SetFormatter(format)
		stream.SetWriter(file)

		if c.StreamRecord {
			stream.WriteFirstLine()

			// To ensure things are flushed correctly
			defer file.Sync()
		}

		if c.Record {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)
			go func() {
				// As currently designed SIGHUP will record only, without stopping the
				// server. SIGINT will record and exit.
				sig := <-sigChan

				store := records.GetStore()

				log.Printf("on=server.Start record-target='%s' record-count=%d request-count=%d\n",
					c.RecordTarget, store.Size(), store.Iterations())

				stream.WriteAll(store)

				if sig == syscall.SIGINT {
					os.Exit(0)
				}
			}()
		}
	}

	return multiListenAndServe(nil)
}
