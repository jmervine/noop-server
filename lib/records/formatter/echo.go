// # formatter/echo
// Formats for the the http response, when in Echo mode.
package formatter

import (
	"fmt"
	"net/http"

	"github.com/jmervine/noop-server/lib/records"
)

const ECHO_TEMPLATE = "status='%d %s' method=%s path=%s headers='%v'"

type Echo struct{}

func (f *Echo) FormatRecordMap(mapped *records.RecordMap) string {
	return commonFormatRecordMap(f, mapped)
}

func (f *Echo) FormatRecord(r records.Record) string {
	path := commonPath(r.Endpoint)

	return fmt.Sprintf(
		ECHO_TEMPLATE,
		r.Status,
		http.StatusText(r.Status),
		r.Method,
		path,
		f.FormatHeader(r.Headers),
	)
}

func (f *Echo) FormatHeader(headers *http.Header) string {
	if headers == nil {
		return ""
	}

	return commonFormatHeader(headers)
}
