package formatter

import (
	"testing"
	"time"

	"github.com/jmervine/noop-server/lib/records"
)

func TestYaml_FormatRecord(t *testing.T) {
	now := time.Now()
	nf := formattedNow(now)
	rec := record()
	rec.Timestamp = now
	format := Yaml{}
	result := format.FormatRecord(rec)

	expect := `timestamp: "` + nf + `"
iterations: 1
endpoint: http://test.host/testing
method: GET
status: 200
sleep: 0
echo: false
headers: {"Foo":["bar"]}
`

	if result != expect {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expect, result)
	}
}

func BenchmarkYaml_FormatRecord(b *testing.B) {
	benchmarkRecordFor(b, Yaml{})
}

func TestYaml_FormatRecordMap(t *testing.T) {
	now := time.Now()
	nf := formattedNow(now)
	format := Yaml{}
	rec := record()
	rec.Timestamp = now
	rmap := records.NewRecordMap()
	rmap.Add(rec)

	result := format.FormatRecordMap(rmap)
	expect := `# Started: ` + now.Format("Mon Jan 2 15:04:05 MST 2006") + `
---
- timestamp: "` + nf + `"
  iterations: 1
  endpoint: http://test.host/testing
  method: GET
  status: 200
  sleep: 0
  echo: false
  headers: {"Foo":["bar"]}
`

	if result != expect {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expect, result)
	}
}

func BenchmarkYaml_FormatRecordMap(b *testing.B) {
	benchmarkRecordMapFor(b, Yaml{})
}
