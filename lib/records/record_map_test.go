package records

import "testing"

func TestRecordMap_Add(t *testing.T) {
	m := RecordMap{}
	r1 := fullRecord()

	m.Add(r1)

	if m.Size() != 1 {
		t.Error("Expected item to be added to RecordMap")
	}

	rec1, _ := m.get(r1.hash())
	if rec1.Iterations != 1 {
		t.Error("Expected Iterations=1, was", rec1.Iterations)
	}

	r2 := fullRecord()
	m.Add(r2)

	rec2, _ := m.get(r2.hash())
	if rec2.Iterations != 2 {
		t.Error("Expected Iterations=2, was", rec2.Iterations)
	}
}
