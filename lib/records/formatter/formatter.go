package formatter

import (
	"net/http"

	"github.com/jmervine/noop-server/lib/records"
)

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
