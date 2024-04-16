package records

import (
	"fmt"
	"sync"
)

var records = new(RecordMap)

type RecordMap struct {
	sync.Map
}

func (rm *RecordMap) Add(rec Record) {
	hash := rec.hash()

	if mapped, ok := rm.LoadOrStore(hash, rec); ok {
		found := mapped.(Record)
		fmt.Printf("Adding: %+v\n", found)
		found.Iterations++

		// Add modified record back to store
		// TODO: Figure out how to store points in RecordMap
		// - Ideally, we'd store pointers and modify in place, but I was having
		// - issues getting that working.
		rm.Store(hash, found)
	}

}

// This function is only used for testing, making private.
// It should be tested in both fetch and read and fetch and
// write conditions before being made public.
func (rm *RecordMap) get(hash string) (Record, bool) {
	if rec, ok := rm.Load(hash); ok {
		return rec.(Record), true
	}
	return Record{}, false
}

// Could be slow and only used for testing on on exit.
func (rm *RecordMap) Size() int {
	var i int
	rm.Range(func(_, _ interface{}) bool {
		i++
		return true
	})
	return i
}
