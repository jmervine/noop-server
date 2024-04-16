package server

import (
	"fmt"
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

	if record {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP)
		go func() {
			<-sigChan

			r := recorder.SerialRecorder{}

			file, err := os.Create("record.txt")
			if err != nil {
				fmt.Printf("Error creating record.txt: %v\n", err)
				return
			}
			defer file.Close()

			r.SetFormatter(format)
			r.SetWriter(file)
			r.WriteAll(records.GetStore())
		}()
	}

	if c.TLSEnabled() {
		return listenAndServeWithTls(svr, c.CertPrivatePath, c.CertKeyPath)
	}

	return svr.ListenAndServe()
}
