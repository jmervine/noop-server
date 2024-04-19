package formatter

import (
	"testing"

	"github.com/jmervine/noop-server/lib/records"
)

func TestJson_FormatRecord(t *testing.T) {
	rec := record()
	format := Json{}
	result := format.FormatRecord(rec)

	expect := "{\"iterations\":1,\"endpoint\":\"http://test.host/testing\",\"method\":\"GET\",\"status\":200,\"sleep\":0,\"echo\":false,\"headers\":{\"Foo\":[\"bar\"]}}"

	if result != expect {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expect, result)
	}
}

func BenchmarkJson_FormatRecord(b *testing.B) {
	rec := record()
	format := Json{}
	for n := 0; n < b.N; n++ {
		format.FormatRecord(rec)
	}
}

func TestJson_FormatRecordMap(t *testing.T) {
	format := Json{}
	rec := record()
	rmap := records.NewRecordMap()
	rmap.Add(rec)

	result := format.FormatRecordMap(rmap)
	expect := "[{\"iterations\":1,\"endpoint\":\"http://test.host/testing\",\"method\":\"GET\",\"status\":200,\"sleep\":0,\"echo\":false,\"headers\":{\"Foo\":[\"bar\"]}}]"

	if result != expect {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expect, result)
	}

	rec = record()
	rmap.Add(rec)
	result = format.FormatRecordMap(rmap)
	expect = "[{\"iterations\":2,\"endpoint\":\"http://test.host/testing\",\"method\":\"GET\",\"status\":200,\"sleep\":0,\"echo\":false,\"headers\":{\"Foo\":[\"bar\"]}}]"

	if result != expect {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expect, result)
	}

	rec = record()
	rec.Status = 300
	rmap.Add(rec)
	result = format.FormatRecordMap(rmap)
	expect = "[{\"iterations\":2,\"endpoint\":\"http://test.host/testing\",\"method\":\"GET\",\"status\":200,\"sleep\":0,\"echo\":false,\"headers\":{\"Foo\":[\"bar\"]}},"
	expect += "{\"iterations\":1,\"endpoint\":\"http://test.host/testing\",\"method\":\"GET\",\"status\":300,\"sleep\":0,\"echo\":false,\"headers\":{\"Foo\":[\"bar\"]}}]"

	if result != expect {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expect, result)
	}

	// TODO: Add datestamp to records.Record for ordering when printing.
	// To avoid sorting issues (for now) I'm going to test pretty formatting with one item, which fails to test parts
	// of this.
	rmap = records.NewRecordMap()
	rmap.Add(rec)
	format.Newline = true
	result = format.FormatRecordMap(rmap)
	expect = `[
	{"iterations":1,"endpoint":"http://test.host/testing","method":"GET","status":300,"sleep":0,"echo":false,"headers":{"Foo":["bar"]}}
]`

	if result != expect {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expect, result)
	}
}

func BenchmarkJson_FormatRecordMap(b *testing.B) {
	rec := record()
	format := Json{}
	rmap := records.NewRecordMap()
	rmap.Add(rec)

	rec = record()
	rmap.Add(rec)

	rec = record()
	rec.Status = 300
	rmap.Add(rec)

	for n := 0; n < b.N; n++ {
		format.FormatRecordMap(rmap)
	}
}
