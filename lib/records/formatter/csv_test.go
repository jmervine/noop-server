package formatter

import "testing"

func BenchmarkCsv(b *testing.B) {
	// This is using Echo to bench commonFormatRecordMap
	benchmarkRecordMapFor(b, Csv{})
}

func BenchmarkCsv_FormatRecord(b *testing.B) {
	benchmarkRecordFor(b, Csv{})
}
