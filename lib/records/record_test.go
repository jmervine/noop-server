package records

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func emptyRecord() Record {
	return Record{
		Iterations: 1,
		Status:     DEFAULT_STATUS,
		Sleep:      0,
		Headers:    &http.Header{},
	}
}

func fullRecord() Record {
	r := emptyRecord()

	e, _ := url.Parse("http://www.example.com")
	r.endpoint = e
	r.Method = "GET"

	return r
}

func TestRecord_parseStatus(t *testing.T) {
	r := Record{}
	r.parseStatus("308")
	if r.Status != http.StatusPermanentRedirect {
		t.Error("Expected status to be parsed and inserted correctly")
	}

	r.parseStatus("ACK")
	if DEFAULT_STATUS != r.Status {
		t.Error("Expected status to be default: ", DEFAULT_STATUS)
	}
}

func TestRecord_parseValuesFromHeader(t *testing.T) {
	// New Record with defaults
	r1 := emptyRecord()

	// Without values in header
	r1.Headers.Add(RECORD_HEADER, "status:200;sleep:10;echo")
	r1.parseValuesFromHeader()
	if DEFAULT_STATUS != r1.Status {
		t.Errorf("Expected status=%v, was %v", DEFAULT_STATUS, r1.Status)
	}

	r2 := emptyRecord()
	s := "status:404;echo;sleep:1s"
	r2.Headers.Add(RECORD_HEADER, s)

	r2.parseValuesFromHeader()
	if http.StatusNotFound != r2.Status {
		t.Errorf("Expected status=%v, was %v", http.StatusNotFound, r2.Status)
	}

	if time.Second != r2.Sleep {
		t.Errorf("Expected sleep=%v, was %v", time.Second, r2.Sleep)
	}

	if !r2.Echo {
		t.Error("Expected echo to be true")
	}
}

func BenchmarkRecord_parseValuesFromHeader(b *testing.B) {
	r := emptyRecord()
	r.Headers.Add("foo", "bar")

	for n := 0; n < b.N; n++ {
		r.parseValuesFromHeader()
	}
}

func TestRecord_parseSleep(t *testing.T) {
	r := Record{}
	r.parseSleep("1s")
	if time.Second != r.Sleep {
		t.Error("Expected sleep to be parsed and inserted correctly from duration string")
	}

	r.parseSleep("1000")
	if time.Second != r.Sleep {
		t.Error("Expected sleep to be parsed and inserted correctly from int string")
	}

	r.parseSleep("ACK")
	if r.Sleep != 0 {
		t.Error("Expected sleep to be default")
	}
}

func TestRecord_hash(t *testing.T) {
	e, _ := url.Parse("http://localhost/foo1")
	r1 := Record{Status: 200, endpoint: e}
	r2 := Record{Status: 200, endpoint: e}
	r3 := Record{Status: 300, endpoint: e}

	if r1.hash() != r2.hash() {
		t.Error("Expected r1 and r1 to have the same hash")
	}

	if r1.hash() == r3.hash() {
		t.Error("Expected r1 and r3 to have difference hashes")
	}
}

func BenchmarkRecord_hash(b *testing.B) {
	for n := 0; n < b.N; n++ {
		rec := fullRecord()
		_ = rec.hash()
	}
}
