package recorder

import (
	"os"
	"sync"

	"github.com/jmervine/noop-server/lib/records"
)

type AsyncRecorder struct {
	SerialRecorder

	// TODO: Does this need to be 'RWMutex', or is 'Mutex' okay?
	mux sync.RWMutex
}

func (r *AsyncRecorder) SetWriter(h *os.File) {
	r.writer = h
}

func (r *AsyncRecorder) WriteOne(rec records.Record) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	str := r.formatter.FormatRecord(&rec)
	if _, err := r.writer.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}

func (r *AsyncRecorder) WriteAll(rec *records.RecordMap) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	str := r.formatter.FormatRecordMap(rec)
	if _, err := r.writer.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}
