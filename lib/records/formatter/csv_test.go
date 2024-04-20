package formatter

import "testing"

func BenchmarkCsv(b *testing.B) {
	// This is using Echo to bench commonFormatRecordMap
	benchmarkRecordMapFor(b, Csv{})
}
