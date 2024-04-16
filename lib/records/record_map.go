package records

import (
	"sync"
)

var defaultStore = NewRecordMap()

// WARNING: This will grow with variable requests
type RecordMap struct {
	mux   sync.RWMutex
	store map[string]Record
}

func NewRecordMap() *RecordMap {
	rm := RecordMap{store: make(map[string]Record)}
	return &rm
}

func (rm *RecordMap) Add(rec Record) {
	hash := rec.hash()

	unlock := rm.rwLocker()
	defer unlock()

	if r, ok := rm.store[hash]; ok {
		r.Iterations++
		rm.store[hash] = r
		return
	}

	rm.store[hash] = rec
}

// Returns a snapshot of k/v pair at the moment
func (rm *RecordMap) Snapshot() map[string]Record {
	unlock := rm.rLocker()
	defer unlock()
	return rm.store
}

// TODO Consider adding RecordMap#Each, currently no use for it

// This function is only used for testing, making private.
// It should be tested in both fetch and read and fetch and
// write conditions before being made public.
func (rm *RecordMap) get(hash string) (Record, bool) {
	unlock := rm.rLocker()
	defer unlock()
	r, ok := rm.store[hash]
	return r, ok
}

func (rm *RecordMap) Size() int {
	unlock := rm.rLocker()
	defer unlock()
	return len(rm.store)
}

func (rm *RecordMap) rwLocker() func() {
	rm.mux.Lock()
	return rm.mux.Unlock
}

func (rm *RecordMap) rLocker() func() {
	rm.mux.RLock()
	return rm.mux.RUnlock
}
