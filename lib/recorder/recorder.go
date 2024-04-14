package recorder

import "net/http"

type RecordFormatter interface {
	Format(*RecordMap) string
	String(*Record) string
	FormatHeader(*http.Header) string
}
