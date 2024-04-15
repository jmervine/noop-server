package recorder

import (
	"os"

	"github.com/jmervine/noop-server/lib/records"
)

type SerialRecorder struct {
	StdRecorder
}

// func (r *SerialRecorder) SetFormatter(f formatter.RecordsFormatter) {
// 	r.formatter = f
// }

func (r *SerialRecorder) SetWriter(h *os.File) {
	r.writer = h
}

func (r *SerialRecorder) WriteOne(rec records.Record) error {
	str := r.formatter.FormatRecord(rec)
	if _, err := r.writer.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}

func (r *SerialRecorder) WriteAll(rec *records.RecordMap) error {
	str := r.formatter.FormatRecordMap(rec)
	if _, err := r.writer.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}
