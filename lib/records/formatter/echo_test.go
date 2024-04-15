package formatter

import (
	"fmt"
	"strings"
	"testing"
)

func Echo_FormatRecordMap(t *testing.T) {
	m := recordMap()
	f := Echo{}
	r := f.FormatRecordMap(m)
	l1 := fmt.Sprintf(ECHO_TEMPLATE, 200, "OK", "GET", "/testing", "Foo:bar")
	l2 := fmt.Sprintf(ECHO_TEMPLATE, 200, "OK", "POST", "/testing", "Foo:bar")
	e := strings.Join([]string{l1, l2}, DEFAULT_HEADER_JOIN)

	if e != r {
		t.Error("Expected format was not returned, got:", r)
	}
}

func Echo_FormatRecord(t *testing.T) {
	f := Echo{}
	rec := record()
	r := f.FormatRecord(rec)
	e := fmt.Sprintf(ECHO_TEMPLATE, 200, "OK", "GET", "/testing", "Foo:bar")

	if r != e {
		t.Errorf("Expected '%s', got '%s', from '%v'", e, r, rec)
	}
}

func Echo_FormatHeaderTest(t *testing.T) {
	t.Skip("This uses a common formatter from 'formatter' which is already tested.")
}
