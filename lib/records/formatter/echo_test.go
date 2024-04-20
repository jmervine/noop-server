package formatter

import (
	"fmt"
	"testing"
)

func TestEcho_FormatRecord(t *testing.T) {
	f := Echo{}
	rec := record()
	r := f.FormatRecord(rec)
	e := fmt.Sprintf(ECHO_TEMPLATE, 200, "OK", "GET", "/testing", "Foo:bar")

	if r != e {
		t.Errorf("Expected '%s', got '%s', from '%v'", e, r, rec)
	}
}

func TestEcho_FormatHeaderTest(t *testing.T) {
	t.Skip("This uses a common formatter from 'formatter' which is already tested.")
}
