package recorder

// Records (REH-kawrds) a record (REH-kuhrd) somewhere.

import (
	"bytes"
	"io"

	"github.com/jmervine/noop-server/lib/records"
	"github.com/jmervine/noop-server/lib/records/formatter"
)

type Recorder interface {
	SetFormatter(formatter.RecordsFormatter)
	SetWriter(*io.Writer)
	WriteOne(records.Record) error
	WriteAll(*records.RecordMap) error
}

type StdRecorder struct {
	formatter formatter.RecordsFormatter
	writer    io.Writer
}

func (r *StdRecorder) SetFormatter(f formatter.RecordsFormatter) {
	r.formatter = f
}

func (r *StdRecorder) SetWriter(h *bytes.Buffer) {
	r.writer = h
}

func (r *StdRecorder) WriteOne(rec records.Record) error {
	str := r.formatter.FormatRecord(rec)
	if _, err := r.writer.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}

func (r *StdRecorder) WriteAll(rec *records.RecordMap) error {
	str := r.formatter.FormatRecordMap(rec)
	if _, err := r.writer.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}
