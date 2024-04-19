package formatter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jmervine/noop-server/lib/records"
)

type Json struct {
	Default
	Newline bool
}

type JsonRecord struct {
	Timestamp  string      `json:"timestamp"`
	Iterations int         `json:"iterations"`
	Endpoint   string      `json:"endpoint"`
	Method     string      `json:"method"`
	Status     int         `json:"status"`
	Sleep      int         `json:"sleep"`
	Echo       bool        `json:"echo"`
	Headers    http.Header `json:"headers"`
}

func (f Json) FormatRecordMap(mapped *records.RecordMap) string {
	snapshot := mapped.Snapshot()
	jrecords := []string{}

	for _, rec := range snapshot {
		jrecords = append(jrecords, f.FormatRecord(rec))
	}

	// TODO: Rethink how "Newline" is provided in formatters.
	// Perhaps remove it and force it when I need it, though that makes
	// NoopClient's formatter tricky perhaps.
	// That said, forcing it here and now...
	//if f.Newline {
	for i, jrec := range jrecords {
		jrecords[i] = fmt.Sprintf("  %s", jrec)
	}

	return fmt.Sprintf("[\n%s\n]", strings.Join(jrecords, ",\n"))
	// }

	// return fmt.Sprintf("[%s]", strings.Join(jrecords, ","))
}

func (f Json) FormatRecord(rec records.Record) string {
	jrec := jsonFromRecord(rec)
	b, err := json.Marshal(jrec)
	if err != nil {
		return fmt.Sprintf("{\"error\": \"%v\"}", err)
	}

	return string(b[:])
}

func (f Json) FormatHeader(header *http.Header) string {
	// Unused
	return ""
}

func (f Json) FormatTimestamp() string {
	// Unused
	return ""
}

// Creating this function, so that I can override it during testing.
// I dislike the pattern of writing addition code for testing, but
// other options see heavy handed.
func formattedNow(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z07:00")
}

func jsonFromRecord(rec records.Record) JsonRecord {
	// Use format RFC3339
	timestamp := formattedNow(rec.Timestamp)
	headers := *rec.Headers
	endpoint := rec.Endpoint()

	return JsonRecord{
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
