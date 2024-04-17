package formatter

import (
	"fmt"
	"strings"
	"testing"
)

func TestEcho_FormatRecordMap(t *testing.T) {
	m := recordMap()
	f := Echo{}
	r := f.FormatRecordMap(m)
	l1 := fmt.Sprintf(ECHO_TEMPLATE, 200, "OK", "GET", "/testing", "Foo:bar")
	l2 := fmt.Sprintf(ECHO_TEMPLATE, 200, "OK", "POST", "/testing", "Foo:bar")
	e := strings.Join([]string{l1, l2}, "\n")

	if e != r {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\n", e, r)
	}
}

func BenchmarkRecordMapAddOne(b *testing.B) {
	m := recordMap()
	f := Echo{}

	for n := 0; n < b.N; n++ {
		_ = f.FormatRecordMap(m)
	}
}

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
