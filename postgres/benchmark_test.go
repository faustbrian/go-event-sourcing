package postgres

// BenchmarkEncodeMetadata measures sorting and encoding a fixed three-entry
// map; no connection, transaction or database operation is timed.

import "testing"

func BenchmarkEncodeMetadata(b *testing.B) {
	metadata := map[string]string{
		"source":      "benchmark",
		"environment": "test",
		"region":      "eu-north-1",
	}

	b.ReportAllocs()
	for b.Loop() {
		_ = encodeMetadata(metadata)
	}
}
