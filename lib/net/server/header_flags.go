package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Parsing header flags takes +/- 20ns, as such, I'm caching the results above
// so that only new header flags will be parsed.
//
// TODO: Consider replacing with a real LRU cache.
//
//	This is dangerous; if someone generates a high volume of random headers,
//	it could grow forever.
var headerFlagCache = make(map[string]*headerFlags)

const FLAG_SLEEP_CAP = (15 * time.Second)

type headerFlags struct {
	status int
	sleep  time.Duration
	echo   bool
}

func (f *headerFlags) Status() int {
	// valid status code range
	if f.status < 100 || f.status > 599 {
		return http.StatusOK
	}
	return f.status
}

func (f *headerFlags) Sleep() {
	if f.sleep == 0 {
		return
	}

	dur := f.sleep
	if dur > FLAG_SLEEP_CAP {
		dur = FLAG_SLEEP_CAP
	}

	time.Sleep(dur)
}

func (f *headerFlags) Echo(req *http.Request) string {
	r := fmt.Sprintf("%d %s", f.Status(), http.StatusText(f.Status()))
	if f.echo {
		r = fmt.Sprintf("status='%s' method=%s path=%s headers='%v' flags='%+v'",
			r, req.Method, req.URL.Path, req.Header, f)
	}

	return r
}

func parseHeaderFlags(header string) *headerFlags {
	// Return cached
	if cf, hit := headerFlagCache[header]; hit {
		return cf
	}

	// Create new to avoid passing the same instance around
	hf := new(headerFlags)

	// MUST DEFAULT TO DESIRED VALUES WITHOUT FLAGS PASSED
	hf.status = http.StatusOK

	// Ensure we cache whatever we do
	defer func() {
		headerFlagCache[header] = hf
	}()

	// if empty, return defaults, which will be cached for next time
	if header == "" {
		return hf
	}

	hSep := ";"
	fSep := "="

	for _, flag := range strings.Split(header, hSep) {
		flagKV := strings.SplitN(flag, fSep, 2)

		key := flagKV[0]
		var val string
		if len(flagKV) == 2 {
			val = flagKV[1]
		}

		switch key {
		case "status":
			hf.status = parseIntOr(val, 500)
		case "sleep":
			hf.sleep = parseDuration(val)
		case "echo":
			hf.echo = true
		}
	}

	return hf
}

func parseIntOr(v string, d int) int {
	if i, e := strconv.ParseInt(v, 10, 16); e == nil {
		return int(i)
	}

	return d
}

func parseDuration(v string) time.Duration {
	// Support direct duration format - e.g. 1s 2ms, etc.
	dur, err := time.ParseDuration(v)
	if err != nil {
		if i, e := strconv.ParseInt(v, 10, 16); e == nil {
			dur = time.Duration(int(i)) * time.Millisecond
		}
	}

	return dur
}
