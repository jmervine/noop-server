package formatter

import (
	"fmt"
	"testing"
)

func TestNoopClient_FormatRecordMap(t *testing.T) {
	t.Skip("NoopClient uses 'commonFormatRecordMap' which is tested by the Echo formatter.")
}

func TestNoopClient_FormatRecord(t *testing.T) {
	f := NoopClient{}
	rec := record()
	r := f.FormatRecord(rec)
	e := fmt.Sprintf(NOOP_CLEINT_TEMPLATE, 1, "GET", "http://test.host/testing", "Foo:bar", 0)

	if r != e {
		t.Errorf("Expected '%s', got '%s', from '%+v'", e, r, rec)
	}
}

func BenchmarkNoopClient_FormatRecord(b *testing.B) {
	benchmarkRecordFor(b, NoopClient{})
}

func TestNoopClient_FormatHeaderTest(t *testing.T) {
	t.Skip("This uses a common formatter from 'formatter' which is already tested.")
}
