package counters

import (
	"os"
	"runtime"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/sdfsdhgjkbmnmxc/ip-counter-go/internal/u32"
)

func NewParallelMMapCounter(name string, newSet func(fileSize int) u32.Set) ParallelMMapCounter {
	return ParallelMMapCounter{
		name:   name,
		newSet: newSet,
	}
}

type ParallelMMapCounter struct {
	name   string
	newSet func(fileSize int) u32.Set
}

func (c ParallelMMapCounter) Name() string { return c.name }

type fileRange struct {
	start int
	end   int
}

func (c ParallelMMapCounter) Count(f *os.File) (int, error) {
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}

	size := int(stat.Size())
	if size == 0 {
		return 0, nil
	}

	numWorkers := runtime.NumCPU()

	minSize := avgIPv4size * 10 * numWorkers
	if size < minSize {
		numWorkers = 1
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return 0, err
	}
	defer func() { _ = syscall.Munmap(data) }()

	seen := c.newSet(size)
	if numWorkers == 1 {
		return c.processRange(data, seen)
	}

	ranges := c.getRanges(numWorkers, data)
	var g errgroup.Group
	for i := 0; i < numWorkers; i++ {
		r := ranges[i]
		g.Go(func() error {
			_, err := c.processRange(data[r.start:r.end], seen)
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return 0, err
	}

	return seen.Count(), nil
}

func (c ParallelMMapCounter) getRanges(numWorkers int, data []byte) []fileRange {
	size := len(data)
	chunkSize := size / numWorkers
	ranges := make([]fileRange, numWorkers)
	start := 0
	for i := 0; i < numWorkers; i++ {
		end := start + chunkSize
		if i == numWorkers-1 {
			end = size
		} else {
			for end < size && data[end] != '\n' {
				end++
			}
			if end < size {
				end++
			}
		}
		ranges[i] = fileRange{start: start, end: end}
		start = end
	}
	return ranges
}

func (c ParallelMMapCounter) processRange(data []byte, seen u32.Set) (int, error) {
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			if i > start {
				ip, err := parseIPv4FromBytes(data[start:i])
				if err != nil {
					return 0, wrapInvalidIPError(err)
				}
				seen.Add(ip)
			}
			start = i + 1
		}
	}

	if start < len(data) {
		ip, err := parseIPv4FromBytes(data[start:])
		if err != nil {
			return 0, wrapInvalidIPError(err)
		}
		seen.Add(ip)
	}

	return seen.Count(), nil
}
