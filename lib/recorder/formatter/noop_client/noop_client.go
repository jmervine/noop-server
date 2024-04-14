// # formater/noop_client
//
// See https://github.com/jmervine/noop-client
// - Specifically examples/*.txt
//
// # Format is '{iterations:-1}|{method:-GET}|{endpoint}|{headers:-}
// 6|GET|http://localhost:3000/request1|User-Agent:noop-client,X-Test:run1
package noop_client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jmervine/noop-server/lib/recorder"
)

const FORMAT_TEMPLATE = "%d|%s|%s|%s"
const FORMAT_HEADER_TEMPLATE = "%s:%s"
const FORMAT_HEADER_JOIN = ";"

type NoopClient struct {
	recorder.RecordFormatter
}

func (f *NoopClient) Format(mapped *recorder.RecordMap) string {
	collect := []string{}
	mapped.Range(func(_, r interface{}) bool {
		collect = append(collect, f.String(r.(*recorder.Record)))
		return true
	})

	return strings.Join(collect, "\n")
}

func (f *NoopClient) String(r *recorder.Record) string {
	return fmt.Sprintf(
		FORMAT_TEMPLATE,
		r.Iterations,
		r.Method,
		r.Endpoint,
		f.FormatHeader(r.Headers),
	)
}

func (f *NoopClient) FormatHeader(headers http.Header) string {
	collect := []string{}
	for key, value := range headers {
		values := strings.Join(value, ",")
		collect = append(collect, fmt.Sprintf(FORMAT_HEADER_TEMPLATE, key, values))
	}

	return strings.Join(collect, FORMAT_HEADER_JOIN)
}
