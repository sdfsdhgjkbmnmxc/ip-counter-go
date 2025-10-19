package counters

import (
	"os"
)

type AutoCounter struct {
	SmallFiles string
	LargeFiles string
}

func (c AutoCounter) Name() string { return "Auto" }

func (c AutoCounter) Count(f *os.File) (int, error) {
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}

	size := stat.Size()
	if size == 0 {
		return 0, nil
	}

	if size/avgIPv4size < 28_000_000 { // see u32_test.go:TestMemoryCrossover
		return Registry.Get(c.SmallFiles).Count(f)
	}

	return Registry.Get(c.LargeFiles).Count(f)
}
