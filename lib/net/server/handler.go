package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const FLAG_HEADER = "X-NoopServerFlags"

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()
	status := http.StatusOK

	defer func() {
		logPrefix := fmt.Sprintf("on=server.handlerFunc method=%s path=%s", r.Method, r.URL.Path)
		log.Printf("%s status=%d took=%v\n", logPrefix, status, time.Since(begin))

		if verbose {
			log.Printf("%s headers:\n%s", logPrefix, r.Header)

			body := &bytes.Buffer{}
			if _, err := io.Copy(body, r.Body); err == nil {
				log.Printf("%s body:\n%s", logPrefix, body.String())
			}
		}
	}()

	var parsed map[string]interface{}
	if header := r.Header.Get(FLAG_HEADER); header != "" {
		parsed = parseHeaderFlags(header)
	}

	if _, ok := parsed["status"]; ok {
		status = parsed["status"].(int)
	}

	// handle parsed flags
	response := fmt.Sprintf("%d %s", status, http.StatusText(status))
	if _, ok := parsed["echo"]; ok {
		response = fmt.Sprintf("method=%s path=%s headers='%v' flags='%v'",
			r.Method, r.URL.Path, r.Header, parsed)
	}

	if dur, ok := parsed["sleep"]; ok {
		time.Sleep(dur.(time.Duration))
	}

	http.Error(w, response, status)
}

// save processing time for repeat calls by caching the parsed header results
var parsedHeaderCache = make(map[string]map[string]interface{})

// Parsing header flags takes +/- 20ns, as such, I'm caching the results above
// so that only new header flags will be parsed.
func parseHeaderFlags(header string) map[string]interface{} {
	if c, ok := parsedHeaderCache[header]; ok {
		return c
	}

	fsplit := ";"
	vsplit := "="

	results := make(map[string]interface{})
	for _, fpair := range strings.Split(header, fsplit) {
		vpair := strings.SplitN(fpair, vsplit, 2)
		key := vpair[0]

		var val string
		if len(vpair) == 2 {
			val = vpair[1]
		}

		switch key {
		case "status":
			results[key] = parsedStatusFlag(val)
		case "sleep":
			results[key] = parsedSleepFlag(val)
		default:
			results[key] = true
		}
	}

	parsedHeaderCache[header] = results
	return results
}

func parsedStatusFlag(v string) int {
	if i, e := strconv.ParseInt(v, 10, 16); e == nil {
		return int(i)
	}

	return 500
}

func parsedSleepFlag(v string) time.Duration {
	// Support direct duration format - e.g. 1s 2ms, etc.
	dur, err := time.ParseDuration(v)
	if err != nil {
		if i, e := strconv.ParseInt(v, 10, 16); e == nil {
			dur = time.Duration(int(i)) * time.Millisecond
		}
	}

	// for safety cap to 15 seconds
	var cap = 15 * time.Second
	if dur > cap {
		return cap
	}

	return dur
}
