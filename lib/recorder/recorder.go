package recorder

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const RECORD_HEADER = "X-Noopserverflags"
const SPLIT_RECORD_HEADER = ";"
const SPLIT_HEADER_VALUE = ":"
const DEFAULT_STATUS = http.StatusOK

// Used to create a string for hashing a Record
const RECORD_HASH_STRING = "status=%d|method=%s|endpoint=%s|header=%#v|sleep=%v|echo=%v"

var Records = RecordMap{}

type RecordMap map[string]*Record

func (rm RecordMap) Add(rec Record) {
	hash := rec.Hash()

	// Add if it doesn't exist
	mapped, ok := rm[hash]

	if !ok {
		rm[hash] = &rec
		return
	}

	// Increment iterations, if it exists.
	mapped.Iterations += 1
}

// TODO: Add Records.Add()

type Record struct {
	Iterations uint16
	Headers    http.Header
	Endpoint   string
	Method     string

	// TODO: Record - Consider using fetcher methods for Status and Sleep to ensure safty
	Status uint16
	Sleep  time.Duration

	// TODO: Record - Support Body in Record, perhapse instead of Echo
	Echo bool
}

func NewRecord(req *http.Request) Record {
	r := Record{}
	// TODO: defer Records.Add()

	// Because this will parse a single request, the iterations will always be 1
	// this field exists to be a counter as they're added to Records.
	r.Iterations = 1

	// Values from http.Request
	r.Endpoint = req.URL.Path
	r.Method = req.Method
	r.Headers = req.Header

	// Values from http.Header
	r.parseValuesFromHeader()

	return r
}

func (r *Record) parseValuesFromHeader() {
	if r.Headers == nil {
		r.Headers = http.Header{}
		return
	}

	header := r.Headers.Get(RECORD_HEADER)

	if header == "" {
		return
	}

	for _, flag := range strings.Split(header, SPLIT_RECORD_HEADER) {
		kv := strings.SplitN(flag, SPLIT_HEADER_VALUE, 2)

		k := kv[0]
		var v string
		if len(kv) == 2 {
			v = kv[1]
		}

		switch k {
		case "status":
			r.parseStatus(v)
		case "sleep":
			r.parseSleep(v)
		case "echo":
			r.Echo = true
		}
	}

}

// TODO: Record.parseStatus - consider return error
func (r *Record) parseStatus(s string) {
	if i, e := strconv.ParseInt(s, 10, 16); e == nil {
		r.Status = uint16(i)
		return
	}

	// Default
	r.Status = DEFAULT_STATUS
}

// TODO: Record.parseSleep - consider return error
func (r *Record) parseSleep(s string) {
	// Support direct duration format - e.g. 1s 2ms, etc.
	dur, err := time.ParseDuration(s)
	if err == nil {
		r.Sleep = dur
		return
	}

	if i, e := strconv.ParseInt(s, 10, 16); e == nil {
		r.Sleep = time.Duration(int(i)) * time.Millisecond
		return
	}

	// Default
	r.Sleep = 0
}

func (r Record) Hash() string {
	hstr := fmt.Sprintf(RECORD_HASH_STRING,
		r.Status, r.Method, r.Endpoint,
		r.Headers, r.Sleep, r.Echo,
	)

	hash := sha256.Sum256([]byte(hstr))
	return fmt.Sprintf("%x", hash)
}
