package counter

import (
	"fmt"
	"os"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Uint32MmapParallel struct {
	Workers int
}

func (c Uint32MmapParallel) Name() string {
	return fmt.Sprintf("uint32_mmap_parallel(w=%d)", c.Workers)
}

func (c Uint32MmapParallel) Count(f *os.File) (int, error) {
	if c.Workers < 1 {
		return 0, fmt.Errorf("workers must be >= 1")
	}

	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}

	size := int(stat.Size())
	if size == 0 {
		return 0, nil
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return 0, err
	}
	defer func() { _ = syscall.Munmap(data) }()

	total, err := c.processChunksParallel(data, c.splitChunks(data))
	if err != nil {
		return 0, err
	}

	return total.Count(), nil
}

func (c Uint32MmapParallel) splitChunks(data []byte) []chunk {
	size := len(data)
	minSize := avgIPv4size * 10

	actualWorkers := c.Workers
	if size < minSize*actualWorkers {
		actualWorkers = size / minSize
		if actualWorkers < 1 {
			actualWorkers = 1
		}
	}

	chunkSize := size / actualWorkers
	chunks := make([]chunk, actualWorkers)

	start := 0
	for i := 0; i < actualWorkers; i++ {
		end := start + chunkSize

		if i == actualWorkers-1 {
			end = size
		} else {
			for end < size && data[end] != '\n' {
				end++
			}
			if end < size {
				end++
			}
		}

		chunks[i] = chunk{start: start, end: end}
		start = end
	}

	return chunks
}

func (c Uint32MmapParallel) processChunksParallel(data []byte, chunks []chunk) (uint32set, error) {
	total := newIPv4Map(0)

	var mux sync.Mutex
	var g errgroup.Group

	for _, ch := range chunks {
		g.Go(func() error {
			result, err := c.processChunk(data, ch)
			if err != nil {
				return err
			}

			mux.Lock()
			defer mux.Unlock()

			total = total.Union(result)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return total, nil
}

func (c Uint32MmapParallel) processChunk(data []byte, ch chunk) (uint32set, error) {
	seen := newIPv4Map(maxCapacity(ch.size() / avgIPv4size))
	lineStart := ch.start

	for i := ch.start; i < ch.end; i++ {
		if data[i] == '\n' {
			if i > lineStart {
				ip, err := parseIPv4FromBytes(data[lineStart:i])
				if err != nil {
					return nil, fmt.Errorf("invalid IP address: %v", err)
				}
				seen.Add(ip)
			}
			lineStart = i + 1
		}
	}

	if lineStart < ch.end {
		ip, err := parseIPv4FromBytes(data[lineStart:ch.end])
		if err != nil {
			return nil, fmt.Errorf("invalid IP address: %v", err)
		}
		seen.Add(ip)
	}

	return seen, nil
}
