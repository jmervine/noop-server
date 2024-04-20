package formatter

import (
	"testing"
	"time"

	"github.com/jmervine/noop-server/lib/records"
)

func TestJson_FormatRecord(t *testing.T) {
	now := time.Now()
	nf := formattedNow(now)
	rec := record()
	rec.Timestamp = now
	format := Json{}
	result := format.FormatRecord(rec)

	expect := "{\"timestamp\":\"" + nf + "\",\"iterations\":1,\"endpoint\":\"http://test.host/testing\",\"method\":\"GET\",\"status\":200,\"sleep\":0,\"echo\":false,\"headers\":{\"Foo\":[\"bar\"]}}"

	if result != expect {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expect, result)
	}
}

func BenchmarkJson_FormatRecord(b *testing.B) {
	benchmarkRecordFor(b, Json{})
}

func TestJson_FormatRecordMap(t *testing.T) {
	now := time.Now()
	nf := formattedNow(now)
	format := Json{}
	rec := record()
	rec.Timestamp = now
	rmap := records.NewRecordMap()
	rmap.Add(rec)

	result := format.FormatRecordMap(rmap)
	expect := "[\n  {\"timestamp\":\"" + nf + "\",\"iterations\":1,\"endpoint\":\"http://test.host/testing\",\"method\":\"GET\",\"status\":200,\"sleep\":0,\"echo\":false,\"headers\":{\"Foo\":[\"bar\"]}}\n]"

	if result != expect {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expect, result)
	}

	rec = record()
	rec.Timestamp = now
	rmap.Add(rec)
	result = format.FormatRecordMap(rmap)
	expect = "[\n  {\"timestamp\":\"" + nf + "\",\"iterations\":2,\"endpoint\":\"http://test.host/testing\",\"method\":\"GET\",\"status\":200,\"sleep\":0,\"echo\":false,\"headers\":{\"Foo\":[\"bar\"]}}\n]"

	if result != expect {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expect, result)
	}

	// Only testing two here because sorting cannot be relied upon.
	rmap = records.NewRecordMap()
	rec = record()
	rec.Timestamp = now
	rmap.Add(rec)
	result = format.FormatRecordMap(rmap)
	expect = "[\n  {\"timestamp\":\"" + nf + "\",\"iterations\":1,\"endpoint\":\"http://test.host/testing\",\"method\":\"GET\",\"status\":200,\"sleep\":0,\"echo\":false,\"headers\":{\"Foo\":[\"bar\"]}}\n]"

	if result != expect {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expect, result)
	}
}

func BenchmarkJson_FormatRecordMap(b *testing.B) {
	benchmarkRecordMapFor(b, Json{})
}
