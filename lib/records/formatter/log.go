package formatter

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jmervine/noop-server/lib/records"
)

const LOG_TEMPLATE = "on=%s method=%s path=%s status=%d took=%v\n"

type Log struct {
	Default
	verbose  bool
	caller   string
	duration time.Duration
	body     io.ReadCloser
}

func NewLogFormatter(fn string, dur time.Duration, body io.ReadCloser, v bool) Log {
	return Log{
		caller:   fn,
		duration: dur,
		body:     body,
		verbose:  v,
	}
}

func (f Log) FormatRecordMap(mapped *records.RecordMap) string {
	return commonFormatRecordMap(f, mapped)
}

func (f Log) FormatRecord(r records.Record) string {
	path := commonPath(r.Path())

	str := fmt.Sprintf(LOG_TEMPLATE, f.caller, r.Method, path, r.Status, f.duration)

	if f.verbose {
		headers := commonFormatHeader(r.Headers)
		str += fmt.Sprintf("on=%sheaders='%s'\n", f.caller, headers)

		body := &bytes.Buffer{}
		if _, err := io.Copy(body, f.body); err == nil {
			bodyStr := body.String()
			if bodyStr != "" {
				str += fmt.Sprintf("on%sbody='%s'\n", f.caller, bodyStr)
			}
		}
	}

	return str
}

func (f Log) FormatHeader(h *http.Header) string {
	return commonFormatHeader(h)
}
