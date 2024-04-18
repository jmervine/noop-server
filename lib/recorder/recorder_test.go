package recorder

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jmervine/noop-server/lib/records"
	"github.com/jmervine/noop-server/lib/records/formatter"
)

func request() *http.Request {
	header := http.Header{}
	header.Add("foo", "bar")

	u, _ := url.Parse("http://localhost/testing")

	req := new(http.Request)

	req.Method = "GET"
	req.URL = u
	req.Header = header

	return req
}

func record() records.Record {
	req := request()
	r := records.NewRecord(req, "test.host:3333", nil)
	r.Iterations = 1
	r.Status = http.StatusOK
	r.Sleep = 0
	r.Echo = false

	return r
}

func TestStdRecord(t *testing.T) {
	result, err := capture(func(w *os.File) {
		rec := record()
		f := &formatter.Echo{}
		recr := &StdRecorder{}

		recr.SetFormatter(f)
		recr.SetWriter(w)

		if recr.formatter == nil {
			t.Error("Expected formatter to be set")
		}

		if recr.writer == nil {
			t.Error("Expected writer to be set")
		}

		if err := recr.WriteOne(rec); err != nil {
			t.Error(err)
		}
	})

	if err != nil {
		t.Error(err)
	}

	expect := fmt.Sprintf(
		formatter.ECHO_TEMPLATE,
		200, "OK", "GET", "/testing", "Foo:bar",
	)

	if result != expect {
		t.Errorf("Expected '%s', got '%s'", expect, result)
	}
}

func capture(fn func(*os.File)) (res string, err error) {
	rpipe, wpipe, err := os.Pipe()

	if err != nil {
		return
	}

	fn(wpipe)
	err = wpipe.Close()
	if err != nil {
		return
	}

	bres, err := io.ReadAll(rpipe)
	if err != nil {
		return
	}

	return string(bres[:]), nil
}
