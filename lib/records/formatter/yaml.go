// formatter/yaml
// The yaml formatting libs found were slow and this is a perscritive output,
// we're going to just write out what we want instead of using a Marshal fn.
package formatter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jmervine/noop-server/lib/records"
)

type Yaml struct {
	Default
	Newline bool
}

func (f Yaml) FormatRecordMap(mapped *records.RecordMap) string {
	template := `- timestamp: "%s"
  iterations: %d
  endpoint: %s
  method: %s
  status: %d
  sleep: %d
  echo: %v
  headers: %s
`

	snapshot := mapped.Snapshot()
	recs := []string{}

	for _, rec := range snapshot {
		r := newCommonSerializer(rec)
		new := fmt.Sprintf(
			template, r.Timestamp, r.Iterations, r.Endpoint, r.Method,
			r.Status, r.Sleep, r.Echo, f.FormatHeader(&r.Headers),
		)
		recs = append(recs, new)
	}

	return f.FormatFirstLine() + strings.Join(recs, "\n")
}

func (f Yaml) FormatRecord(rec records.Record) string {
	template := `timestamp: "%s"
iterations: %d
endpoint: %s
method: %s
status: %d
sleep: %d
echo: %v
headers: %s
`

	r := newCommonSerializer(rec)
	return fmt.Sprintf(
		template, r.Timestamp, r.Iterations, r.Endpoint, r.Method,
		r.Status, r.Sleep, r.Echo, f.FormatHeader(&r.Headers),
	)
}

func (y Yaml) FormatHeader(headers *http.Header) string {
	// Because json.Marshal is very fast, and this result is more fluid we are
	// going to use it here.
	b, err := json.Marshal(headers)
	if err != nil {
		return fmt.Sprintf("{\"error\": \"when marshaling headers got: %v\",\"headers\":\"%+v\"}\n", err, headers)
	}

	return string(b[:])
}

func (f Yaml) FormatFirstLine() string {
	return fmt.Sprintf("# Started: %s\n---\n", time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
}
