package counter

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkCounter(b *testing.B) {
	testFiles, err := filepath.Glob("../../testdata/*.txt")
	if err != nil {
		b.Fatalf("Failed to find test files: %v", err)
	}
	if len(testFiles) == 0 {
		b.Skip("No test files found in testdata/")
	}

	for _, path := range testFiles {
		filename := filepath.Base(path)
		b.Run(filename, func(b *testing.B) {
			for method, counter := range Counters {
				b.Run(method, func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						f, err := os.Open(path)
						if err != nil {
							b.Fatalf("Failed to open %s: %v", path, err)
						}

						_, err = counter.Count(f)
						f.Close()
						if err != nil {
							b.Fatalf("Count failed: %v", err)
						}
					}
				})
			}
		})
	}
}
