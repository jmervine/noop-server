package formatter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jmervine/noop-server/lib/records"
)

const DEFAULT_HEADER_TEMPLATE = "%s:%s"
const DEFAULT_HEADER_JOIN = ";"

var Formatter recordsFormatter = new(NoopClient)

type recordsFormatter interface {
	FormatRecordMap(*records.RecordMap) string
	FormatRecord(*records.Record) string
	FormatHeader(*http.Header) string
}

// For future use, when there are more formatters available.
func SetFormatter(formatter string) recordsFormatter {
	var rf recordsFormatter
	defer func() {
		Formatter = rf
	}()

	switch formatter {
	case "noop_client":
	case "noopclient":
	default:
		rf = &NoopClient{}
	}

	return rf
}

// Common - reuse in more than one
func commonFormatRecordMap(f recordsFormatter, mapped *records.RecordMap) string {
	collect := []string{}
	mapped.Range(func(_, r interface{}) bool {
		collect = append(collect, f.FormatRecord(r.(*records.Record)))
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
