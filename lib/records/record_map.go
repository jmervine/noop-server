package records

import "sync"

var records = RecordMap{}

type RecordMap struct {
	sync.Map
}

func (rm *RecordMap) Add(rec Record) {
	hash := rec.hash()

	if mapped, ok := rm.LoadOrStore(hash, &rec); ok {
		found := mapped.(*Record)
		found.Iterations++
	}
}

func (rm *RecordMap) Get(hash string) (*Record, bool) {
	if rec, ok := rm.Load(hash); ok {
		return rec.(*Record), true
	}
	return nil, false
}

// Could be slow and currently onle used for testing, so I'm
// making it internal only.
func (rm *RecordMap) size() int {
	var i int
	rm.Range(func(_, _ interface{}) bool {
		i++
		return true
	})
	return i
}
