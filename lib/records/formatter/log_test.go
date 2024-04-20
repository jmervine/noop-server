package formatter

import "testing"

func BenchmarkLog_FormatRecord(b *testing.B) {
	benchmarkRecordFor(b, Log{})
}
