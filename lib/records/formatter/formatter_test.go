package formatter

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/jmervine/noop-server/lib/records"
)

// These functions can be used in all formatter tests.
func recordMap() *records.RecordMap {
	m := records.NewRecordMap()
	r1 := record()
	m.Add(r1)

	r2 := record()
	r2.Method = "POST"
	m.Add(r2)

	return m
}

func request() *http.Request {
	u, _ := url.Parse("http://test.host/testing")
	req := new(http.Request)
	req.Method = "GET"
	req.URL = u
	req.Header = headers(true)

	return req
}

func record() records.Record {
	req := request()
	r := records.NewRecord(req, "")
	r.Iterations = 1
	r.Status = http.StatusOK
	r.Sleep = 0

	return r
}

// Because sorting isn't assured, we have the option to say give me
// only one header for when we're not actualyl testing header formatting.
// Otherwise, we'll want two, to make sure that joining is correct, but
// we'll have to test that it's either order. See example in
// 'TestFormatter_commonFormatHeader' below.
func headers(onlyOne bool) http.Header {
	h := http.Header{}
	h.Add("foo", "bar")
	if !onlyOne {
		h.Add("bah", "bin")
	}
	return h
}

// Skipping "TestFormatter_commonFormatRecordMap()", as it requires a
// formatter and as such, will tested by implemented formatters.
func TestFormatter_commonFormatRecordMap(t *testing.T) {
	t.Skip("Skipping this test, as it requires a formatter and as such, will tested by implemented formatters.")
}

func TestFormatter_commonFormatHeader(t *testing.T) {
	h := headers(false)
	r := commonFormatHeader(&h)

	// Because order cannot be assured, we're going to test for either of
	// these two strings.
	e1 := "Foo:bar;Bah:bin"
	e2 := "Bah:bin;Foo:bar"

	if r != e1 && r != e2 {
		t.Errorf("Expected '%s' or '%s', got '%s', from '%v'", e1, e2, r, h)
	}
}
