package formatter

import (
	"encoding/json"
	"fmt"
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

	for i, jrec := range jrecords {
		jrecords[i] = fmt.Sprintf("  %s", jrec)
	}

	return fmt.Sprintf("[\n%s\n]", strings.Join(jrecords, ",\n"))
}

func (f Json) FormatRecord(rec records.Record) string {
	jrec := newCommonSerializer(rec)
	b, err := json.Marshal(jrec)
	if err != nil {
		return fmt.Sprintf("{\"error\": \"%v\"}\n", err)
	}

	return string(b[:])
}
