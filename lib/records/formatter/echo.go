// # formatter/echo
// Formats for the the http response, when in Echo mode.
package formatter

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jmervine/noop-server/lib/records"
)

const ECHO_TEMPLATE = "status='%d %s' method=%s path=%s headers='%v'"

type Echo struct{}

func (f *Echo) FormatRecordMap(mapped *records.RecordMap) string {
	return commonFormatRecordMap(f, mapped)
}

func (f *Echo) FormatRecord(r records.Record) string {
	var path string
	parsed, err := url.Parse(r.Endpoint)
	if err == nil {
		path = parsed.Path
	}

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
