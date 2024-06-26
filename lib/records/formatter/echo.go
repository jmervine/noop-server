// # formatter/echo
// Formats for the the http response, when in Echo mode.
package formatter

import (
	"fmt"
	"net/http"

	"github.com/jmervine/noop-server/lib/records"
)

const ECHO_TEMPLATE = "status='%d %s' method=%s path=%s headers='%v'"

type Echo struct {
	Default
}

func (f Echo) FormatRecord(r records.Record) string {
	path := commonPath(r.Path())

	return fmt.Sprintf(
		ECHO_TEMPLATE,
		r.Status,
		http.StatusText(r.Status),
		r.Method,
		path,
		f.FormatHeader(r.Headers),
	)
}
