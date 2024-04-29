package records

import (
	"testing"
)

func TestRecordMap_Add(t *testing.T) {
	m := NewRecordMap()
	r1 := fullRecord()

	m.Add(r1)

	size := m.Size()
	if size != 1 {
		t.Error("Expected 1 items to be added to RecordMap, got len", size)
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

	r3 := fullRecord()
	r3.Method = "POST"
	m.Add(r3)
	size = m.Size()
	if size != 2 {
		t.Error("Expected 2 items to be added to RecordMap, got len", size)
	}
}

func BenchmarkRecordMap_AddOne(b *testing.B) {
	rm := GetStore()

	for n := 0; n < b.N; n++ {
		rec := fullRecord()
		rm.Add(rec)
	}
}

func BenchmarkRecordMap_AddMany(b *testing.B) {
	rm := GetStore()

	for n := 0; n < b.N; n++ {
		rec := fullRecord()
		rec.Iterations = n
		rm.Add(rec)
	}
}

func TestRecordMap_Size(t *testing.T) {
	m := NewRecordMap()
	r1 := fullRecord()
	r2 := fullRecord()

	m.Add(r1)
	m.Add(r2)

	if m.Size() != 1 {
		t.Error("Expected 1, got", m.Size())
	}

	r3 := fullRecord()
	r3.Method = "DELETE"
	m.Add(r3)

	if m.Size() != 2 {
		t.Error("Expected 2, got", m.Size())
	}
}

func TestRecordMap_Interations(t *testing.T) {
	m := NewRecordMap()
	r1 := fullRecord()
	r1.Iterations = 5

	r2 := fullRecord()
	r2.Iterations = 5

	m.Add(r1)
	m.Add(r2)

	if m.Iterations() != 10 {
		t.Error("Expected 10, got", m.Iterations())
	}

	r3 := fullRecord()
	r3.Iterations = 5
	r3.Method = "DELETE"
	m.Add(r3)

	if m.Iterations() != 15 {
		t.Error("Expected 15, got", m.Iterations())
	}
}
