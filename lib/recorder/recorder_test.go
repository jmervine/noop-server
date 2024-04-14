package recorder

import (
	"net/http"
	"testing"
	"time"
)

func emptyRecord() Record {
	return Record{
		Iterations: 1,
		Status:     DEFAULT_STATUS,
		Sleep:      0,
		Headers:    http.Header{},
	}
}

func fullRecord() Record {
	r := emptyRecord()
	r.Endpoint = "http://www.example.com"
	r.Method = "GET"
	return r
}

func TestRecordMap_Add(t *testing.T) {
	m := RecordMap{}
	r := fullRecord()

	m.Add(r)

	if m.size() != 1 {
		t.Error("Expected item to be added to RecordMap")
	}

	rec, _ := m.Get(r.hash())
	if rec.Iterations != 1 {
		t.Error("Expected Iterations=1, was", rec.Iterations)
	}

	m.Add(r)

	rec, _ = m.Get(r.hash())
	if rec.Iterations != 2 {
		t.Error("Expected Iterations=2, was", rec.Iterations)
	}
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
	r := emptyRecord()

	// Without values in header
	r.Headers.Add("foo", "bar")
	r.parseValuesFromHeader()
	if DEFAULT_STATUS != r.Status {
		t.Errorf("Expected status=%v, was %v", DEFAULT_STATUS, r.Status)
	}

	s := "status:404;echo;sleep:1s"
	r.Headers.Add(RECORD_HEADER, s)

	r.parseValuesFromHeader()
	if http.StatusNotFound != r.Status {
		t.Errorf("Expected status=%v, was %v", http.StatusNotFound, r.Status)
	}

	if time.Second != r.Sleep {
		t.Errorf("Expected sleep=%v, was %v", time.Second, r.Sleep)
	}

	if !r.Echo {
		t.Error("Expected echo to be true")
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
	r1 := Record{Status: 200, Endpoint: "http://localhost/foo1"}
	r2 := Record{Status: 200, Endpoint: "http://localhost/foo1"}
	r3 := Record{Status: 300, Endpoint: "http://localhost/foo1"}

	if r1.hash() != r2.hash() {
		t.Error("Expected r1 and r1 to have the same hash")
	}

	if r1.hash() == r3.hash() {
		t.Error("Expected r1 and r3 to have difference hashes")
	}
}
