// # formatter
//
// Used to format records
package formatter

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jmervine/noop-server/lib/records"
)

const DEFAULT_HEADER_TEMPLATE = "%s:%s"
const DEFAULT_HEADER_JOIN = ";"

type CommonRecordSerializer struct {
	Timestamp  string      `json:"timestamp" yaml:"timestamp"`
	Iterations int         `json:"iterations" yaml:"iterations"`
	Endpoint   string      `json:"endpoint" yaml:"endpoint"`
	Method     string      `json:"method" yaml:"method"`
	Status     int         `json:"status" yaml:"status"`
	Sleep      int         `json:"sleep" yaml:"sleep"`
	Echo       bool        `json:"echo" yaml:"echo"`
	Headers    http.Header `json:"headers" yaml:"headers,omitempty,flow"`
}

// RecordFormatter interface
//   - RecordMap needs to be passed as a pointer to ensure thread
//     safety.
type RecordsFormatter interface {
	// Only used by noop currently, but making available for interface.
	SetNewline(bool)

	FormatRecordMap(*records.RecordMap) string
	FormatRecord(records.Record) string
	FormatHeader(*http.Header) string
	FormatTimestamp() string
}

func NewFromString(formatter string) RecordsFormatter {
	var rf RecordsFormatter

	switch formatter {
	case "noop_client":
	case "noopclient":
		rf = NoopClient{}
	case "echo":
		rf = Echo{}
	default:
		rf = Default{}
	}

	return rf
}

type Default struct{}

func (f Default) SetNewline(_ bool) {}

func (f Default) FormatRecordMap(mapped *records.RecordMap) string {
	return commonFormatRecordMap(f, mapped)
}

func (f Default) FormatRecord(r records.Record) string {
	return fmt.Sprintf("%d %s", r.Status, http.StatusText(r.Status))
}

func (f Default) FormatHeader(headers *http.Header) string {
	if headers == nil {
		return ""
	}

	return commonFormatHeader(headers)
}

func (f Default) FormatTimestamp() string {
	return "" // default to no timestamp
}

// Common - reuse in more than one
func commonFormatRecordMap(f RecordsFormatter, mapped *records.RecordMap) string {
	records := mapped.Snapshot()
	collect := make([]string, 0, len(records))

	for _, record := range mapped.Snapshot() {
		collect = append(collect, f.FormatRecord(record))
	}

	return strings.Join(collect, "\n")
}

func commonFormatHeader(headers *http.Header) string {
	collect := []string{}
	for key, value := range *headers {
		values := strings.Join(value, ",")
		collect = append(collect, fmt.Sprintf(DEFAULT_HEADER_TEMPLATE, key, values))
	}

	return strings.Join(collect, DEFAULT_HEADER_JOIN)
}

func commonPath(s string) string {
	parsed, err := url.Parse(s)
	if err == nil {
		return parsed.Path
	}
	return ""
}

func newCommonSerializer(rec records.Record) CommonRecordSerializer {
	// Use format RFC3339
	timestamp := formattedNow(rec.Timestamp)
	headers := *rec.Headers
	endpoint := rec.Endpoint()

	return CommonRecordSerializer{
		Timestamp:  timestamp,
		Iterations: rec.Iterations,
		Endpoint:   endpoint,
		Method:     rec.Method,
		Status:     rec.Status,
		Sleep:      int(rec.Sleep * time.Millisecond),
		Echo:       rec.Echo,
		Headers:    headers,
	}
}

func formattedNow(t time.Time) string {
	// Format using RFC3339Nano
	return t.Format("2006-01-02T15:04:05.999999999Z07:00")
}
