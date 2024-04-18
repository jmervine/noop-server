package formatter

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestEcho_FormatRecordMap(t *testing.T) {
	m := recordMap()
	f := Echo{}

	// Remove trailing newline and split on newline, and then sort
	got := strings.Split(strings.TrimRight(f.FormatRecordMap(m), "\n"), "\n")
	sort.Strings(got)

	// because I can't trust order
	exp := []string{
		fmt.Sprintf(ECHO_TEMPLATE, 200, "OK", "GET", "/testing", "Foo:bar"),
		fmt.Sprintf(ECHO_TEMPLATE, 200, "OK", "POST", "/testing", "Foo:bar"),
	}
	sort.Strings(exp)

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\n", exp, got)
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
