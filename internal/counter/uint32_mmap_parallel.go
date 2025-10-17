package counter

import (
	"fmt"
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

func (c chunk) size() int {
	return c.end - c.start
}

func (c Uint32MmapParallel) Name() string {
	return fmt.Sprintf("uint32_mmap_parallel(workers=%d)", c.Workers)
}

func (c Uint32MmapParallel) minChunkSize() int {
	return avgIPv4size * 10
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

	return len(c.mergeResults(c.processChunksParallel(data, c.splitChunks(data)))), nil
}

func (c Uint32MmapParallel) splitChunks(data []byte) []chunk {
	size := len(data)
	minSize := c.minChunkSize()

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

func (c Uint32MmapParallel) processChunksParallel(data []byte, chunks []chunk) []map[uint32]struct{} {
	var wg sync.WaitGroup
	results := make([]map[uint32]struct{}, len(chunks))

	for i, ch := range chunks {
		wg.Add(1)
		go func(idx int, ch chunk) {
			defer wg.Done()
			results[idx] = c.processChunk(data, ch)
		}(i, ch)
	}

	wg.Wait()
	return results
}

func (c Uint32MmapParallel) processChunk(data []byte, ch chunk) map[uint32]struct{} {
	seen := make(map[uint32]struct{}, maxCapacity(ch.size()/avgIPv4size))
	lineStart := ch.start

	for i := ch.start; i < ch.end; i++ {
		if data[i] == '\n' {
			if i > lineStart {
				if ip, err := parseIPv4FromBytes(data[lineStart:i]); err == nil {
					seen[ip] = struct{}{}
				}
			}
			lineStart = i + 1
		}
	}

	if lineStart < ch.end {
		if ip, err := parseIPv4FromBytes(data[lineStart:ch.end]); err == nil {
			seen[ip] = struct{}{}
		}
	}

	return seen
}

func (c Uint32MmapParallel) mergeResults(results []map[uint32]struct{}) map[uint32]struct{} {
	totalSize := 0
	for _, result := range results {
		totalSize += len(result)
	}

	total := make(map[uint32]struct{}, maxCapacity(totalSize))
	for _, result := range results {
		for ip := range result {
			total[ip] = struct{}{}
		}
	}

	return total
}
