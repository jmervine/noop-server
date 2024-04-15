// # formatter
//
// Used to format records
package formatter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jmervine/noop-server/lib/records"
)

const DEFAULT_HEADER_TEMPLATE = "%s:%s"
const DEFAULT_HEADER_JOIN = ";"

type RecordsFormatter interface {
	FormatRecordMap(*records.RecordMap) string
	FormatRecord(records.Record) string
	FormatHeader(*http.Header) string
}

func NewFromString(formatter string) RecordsFormatter {
	var rf RecordsFormatter

	switch formatter {
	case "noop_client":
	case "noopclient":
		rf = &NoopClient{}
	case "echo":
		rf = &Echo{}
	default:
		rf = &Default{}
	}

	return rf
}

type Default struct{}

func (f *Default) FormatRecordMap(mapped *records.RecordMap) string {
	return commonFormatRecordMap(f, mapped)
}

func (f *Default) FormatRecord(r records.Record) string {
	return fmt.Sprintf("%d %s", r.Status, http.StatusText(r.Status))
}

func (f *Default) FormatHeader(headers *http.Header) string {
	if headers == nil {
		return ""
	}

	return commonFormatHeader(headers)
}

// Common - reuse in more than one
func commonFormatRecordMap(f RecordsFormatter, mapped *records.RecordMap) string {
	collect := []string{}
	mapped.Range(func(_, r interface{}) bool {
		collect = append(collect, f.FormatRecord(r.(records.Record)))
		return true
	})

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
