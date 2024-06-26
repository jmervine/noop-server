// # formatter/noop_client
//
// See https://github.com/jmervine/noop-client
// - Specifically examples/*.txt
package formatter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jmervine/noop-server/lib/records"
)

// # Format is '{iterations:-1}|{method:-GET}|{endpoint}|{headers:-}|{sleep}
const NOOP_CLEINT_TEMPLATE = "%d|%s|%s|%s|%v"

type NoopClient struct {
	Default
	Newline bool
}

func (f NoopClient) FormatRecordMap(mapped *records.RecordMap) string {
	return f.FormatFirstLine() + commonFormatRecordMap(f, mapped) + "\n"
}

func (f NoopClient) FormatRecord(r records.Record) string {
	return fmt.Sprintf(
		NOOP_CLEINT_TEMPLATE,
		r.Iterations,
		r.Method,
		r.Endpoint(),
		f.FormatHeader(r.Headers),
		int64(r.Sleep*time.Millisecond),
	)
}

func (f NoopClient) FormatHeader(headers *http.Header) string {
	if headers == nil {
		return ""
	}

	return commonFormatHeader(headers)
}

func (f NoopClient) FormatFirstLine() string {
	return fmt.Sprintf("# Started: %s\n", time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
}
