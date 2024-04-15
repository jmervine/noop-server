package recorder

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/jmervine/noop-server/lib/records"
	"github.com/jmervine/noop-server/lib/records/formatter"
)

func record() records.Record {
	header := &http.Header{}
	header.Add("foo", "bar")

	return records.Record{
		Iterations: 1,
		Headers:    header,
		Endpoint:   "http://localhost/testing",
		Method:     "GET",
		Status:     http.StatusOK,
		Sleep:      0,
		Echo:       false,
	}
}

func TestStdRecord(t *testing.T) {
	var buf bytes.Buffer

	rec := record()
	f := &formatter.Echo{}
	r := StdRecorder{}

	r.SetFormatter(f)
	r.SetWriter(&buf)
	r.WriteOne(rec)

	result := buf.String()
	expect := fmt.Sprintf(
		formatter.ECHO_TEMPLATE,
		200, "OK", "GET", "/testing", "Foo:bar",
	)

	if result != expect {
		t.Errorf("Expected '%s', got '%s'", expect, result)
	}
}
