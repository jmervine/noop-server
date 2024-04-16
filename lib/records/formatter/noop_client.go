// # formatter/noop_client
//
// See https://github.com/jmervine/noop-client
// - Specifically examples/*.txt
//
// # Format is '{iterations:-1}|{method:-GET}|{endpoint}|{headers:-}
// 6|GET|http://localhost:3000/request1|User-Agent:noop-client,X-Test:run1
package formatter

import (
	"fmt"
	"net/http"

	"github.com/jmervine/noop-server/lib/records"
)

const NOOP_CLEINT_TEMPLATE = "%d|%s|%s|%s"

type NoopClient struct{}

func (f NoopClient) FormatRecordMap(mapped *records.RecordMap) string {
	return commonFormatRecordMap(f, mapped)
}

func (f NoopClient) FormatRecord(r *records.Record) string {
	return fmt.Sprintf(
		NOOP_CLEINT_TEMPLATE,
		r.Iterations,
		r.Method,
		r.Endpoint,
		f.FormatHeader(r.Headers),
	)
}

func (f NoopClient) FormatHeader(headers *http.Header) string {
	if headers == nil {
		return ""
	}

	return commonFormatHeader(headers)
}
