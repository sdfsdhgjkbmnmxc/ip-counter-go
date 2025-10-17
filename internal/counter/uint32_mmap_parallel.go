package counter

import (
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
)

type Uint32MmapParallel struct {
	Workers int
}

type chunk struct {
	start, end int
}

func (c Uint32MmapParallel) Name() string {
	return fmt.Sprintf("uint32_mmap_parallel_%d", c.Workers)
}

func (c Uint32MmapParallel) Count(r io.Reader) (int, error) {
	if c.Workers < 1 {
		return 0, fmt.Errorf("workers must be >= 1")
	}

	file, ok := r.(*os.File)
	if !ok {
		return 0, fmt.Errorf("mmap requires *os.File")
	}

	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	size := int(stat.Size())
	if size == 0 {
		return 0, nil
	}

	data, err := syscall.Mmap(int(file.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return 0, err
	}
	defer func() { _ = syscall.Munmap(data) }()

	chunks := c.splitChunks(data)
	results := c.processChunksParallel(data, chunks)
	total := c.mergeResults(results)

	return len(total), nil
}

func (c Uint32MmapParallel) splitChunks(data []byte) []chunk {
	size := len(data)
	chunkSize := size / c.Workers
	chunks := make([]chunk, c.Workers)

	start := 0
	for i := 0; i < c.Workers; i++ {
		end := start + chunkSize

		if i == c.Workers-1 {
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

func (c Uint32MmapParallel) processChunksParallel(data []byte, chunks []chunk) []map[uint32]struct{} {
	var wg sync.WaitGroup
	results := make([]map[uint32]struct{}, len(chunks))

	for i, ch := range chunks {
		wg.Add(1)
		go func(idx int, start, end int) {
			defer wg.Done()
			results[idx] = c.processChunk(data, start, end)
		}(i, ch.start, ch.end)
	}

	wg.Wait()
	return results
}

func (c Uint32MmapParallel) processChunk(data []byte, start, end int) map[uint32]struct{} {
	seen := make(map[uint32]struct{})
	lineStart := start

	for i := start; i < end; i++ {
		if data[i] == '\n' {
			if i > lineStart {
				if ip, err := parseIPv4FromBytes(data[lineStart:i]); err == nil {
					seen[ip] = struct{}{}
				}
			}
			lineStart = i + 1
		}
	}

	if lineStart < end {
		if ip, err := parseIPv4FromBytes(data[lineStart:end]); err == nil {
			seen[ip] = struct{}{}
		}
	}

	return seen
}

func (c Uint32MmapParallel) mergeResults(results []map[uint32]struct{}) map[uint32]struct{} {
	total := make(map[uint32]struct{})
	for _, result := range results {
		for ip := range result {
			total[ip] = struct{}{}
		}
	}
	return total
}
