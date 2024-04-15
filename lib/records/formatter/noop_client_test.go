package formatter

import (
	"fmt"
	"testing"
)

func NoopClient_FormatRecordMap(t *testing.T) {
	t.Skip("NoopClient uses 'commonFormatRecordMap' which is tested by the Echo formatter.")
}

func NoopClient_FormatRecord(t *testing.T) {
	f := NoopClient{}
	rec := record()
	r := f.FormatRecord(&rec)
	e := fmt.Sprintf(NOOP_CLEINT_TEMPLATE, 1, "GET", "http://localhost/testing", "Foo:bar")

	if r != e {
		t.Errorf("Expected '%s', got '%s', from '%v'", e, r, rec)
	}
}

func NoopClient_FormatHeaderTest(t *testing.T) {
	t.Skip("This uses a common formatter from 'formatter' which is already tested.")
}
