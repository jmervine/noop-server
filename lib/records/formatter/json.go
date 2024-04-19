package formatter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jmervine/noop-server/lib/records"
)

type Json struct {
	Default
	Newline bool
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
	jrec := newCommonSerializer(rec)
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
