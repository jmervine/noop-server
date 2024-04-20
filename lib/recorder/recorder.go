package recorder

// Records (REH-kawrds) a record (REH-kuhrd) somewhere.

import (
	"os"

	"github.com/jmervine/noop-server/lib/records"
	"github.com/jmervine/noop-server/lib/records/formatter"
	"github.com/pkg/errors"
)

// TIL os.Pipe can be used if you need buffer like functionality
type Recorder interface {
	SetFormatter(formatter.RecordsFormatter)
	SetWriter(*os.File)
	WriteOne(records.Record) error
	WriteAll(*records.RecordMap) error
	WriteString(string) error
	WriteFirstLine()
}

// This will support anything that implements the 'io.Writer' interface.
type StdRecorder struct {
	formatter formatter.RecordsFormatter
	writer    *os.File
}

// This is idempotent, and will fail quietly.
func (r *StdRecorder) WriteFirstLine() {
	ts := r.formatter.FormatFirstLine()
	if ts != "" {
		r.writer.Write([]byte(ts))
		r.writer.Sync()
	}
}

func (r *StdRecorder) SetFormatter(f formatter.RecordsFormatter) {
	r.formatter = f
}

func (r *StdRecorder) SetWriter(f *os.File) {
	r.writer = f
}

func (r *StdRecorder) WriteOne(rec records.Record) error {
	if r.writer == nil {
		return errors.Errorf("writer is not set in: %#v", r)
	}

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

func (r *StdRecorder) WriteString(str string) error {
	if _, err := r.writer.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}
