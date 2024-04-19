package formatter

import (
	"fmt"
	"net/http"

	"github.com/jmervine/noop-server/lib/records"
	"gopkg.in/yaml.v2"
)

type Yaml struct {
	Default
	Newline bool
}

func (f Yaml) FormatRecordMap(mapped *records.RecordMap) string {
	snapshot := mapped.Snapshot()
	recs := []CommonRecordSerializer{}

	for _, rec := range snapshot {
		recs = append(recs, newCommonSerializer(rec))
	}

	b, err := yaml.Marshal(recs)
	if err != nil {
		return fmt.Sprintf("---\nerror: %s\n", err)
	}

	return string(b[:])
}

func (f Yaml) FormatRecord(rec records.Record) string {
	r := newCommonSerializer(rec)
	b, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Sprintf("---\nerror: %s\n", err)
	}

	return string(b[:])
}

func (f Yaml) FormatHeader(header *http.Header) string {
	// Unused
	return ""
}

func (f Yaml) FormatTimestamp() string {
	// Unused
	return ""
}
