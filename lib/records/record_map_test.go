package records

import "testing"

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
