package counter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func BenchmarkCounter(b *testing.B) {
	testFiles, err := filepath.Glob("../../testdata/*.txt")
	if err != nil {
		b.Fatalf("Failed to find test files: %v", err)
	}

	var sizes = make(map[string]int64, len(testFiles))
	for _, path := range testFiles {
		fi, err := os.Stat(path)
		if err != nil {
			b.Fatalf("Failed to stat file %s: %v", path, err)
		}
		sizes[path] = fi.Size()
	}

	sort.Slice(testFiles, func(i, j int) bool {
		return sizes[testFiles[i]] < sizes[testFiles[j]]
	})

	for _, path := range testFiles {
		b.Run(filepath.Base(path), func(b *testing.B) {
			for _, counter := range Counters {
				b.Run(counter.Name(), func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						err := func() error {
							f, err := os.Open(path)
							if err != nil {
								return fmt.Errorf("failed to open file %s: %w", path, err)
							}
							defer func() { _ = f.Close() }()

							res, err := counter.Count(f)
							if err != nil {
								return fmt.Errorf("counter.Count() failed: %w", err)
							}

							if res == 0 {
								return errors.New("counter.Count() returned 0")
							}

							return nil
						}()
						if err != nil {
							b.Fatalf("Error during counting: %v", err)
						}
					}
				})
			}
		})
	}
}
