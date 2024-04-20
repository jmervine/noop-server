// # formatter/csv
//
// See https://github.com/jmervine/noop-client
// - Specifically examples/*.txt
package formatter

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jmervine/noop-server/lib/records"
)

// # Format is time,iter,endpoint,method,status,sleep,echo,headers
const CSV_HEADER = "timestamp,iterations,method,endpoint,status,sleep,headers\n"
const CSV_ROW_TEMPLATE = "%s,%d,%s,%s,%d,%d,%s\n"

type Csv struct {
	Default
	Newline bool
}

func (f Csv) FormatRecordMap(mapped *records.RecordMap) string {
	collect := []string{}
	for _, record := range mapped.Snapshot() {
		collect = append(collect, f.FormatRecord(record))
	}
	return strings.Join(collect, "")
}

func (f Csv) FormatRecord(r records.Record) string {
	out := fmt.Sprintf(
		CSV_ROW_TEMPLATE,
		formattedNow(r.Timestamp),
		r.Iterations,
		r.Method,
		r.Endpoint(),
		r.Status,
		int64(r.Sleep*time.Millisecond),
		f.FormatHeader(r.Headers),
	)

	return out
}

func (f Csv) FormatHeader(headers *http.Header) string {
	if headers == nil {
		return ""
	}

	return "\"" + commonFormatHeader(headers) + "\""
}

func (f Csv) FormatFirstLine() string {
	return CSV_HEADER
}
