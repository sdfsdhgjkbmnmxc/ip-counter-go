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

	return len(c.processChunksParallel(data, c.splitChunks(data))), nil
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

func (c Uint32MmapParallel) processChunksParallel(data []byte, chunks []chunk) IPv4set {
	resultChan := make(chan IPv4set, len(chunks))

	var wg sync.WaitGroup
	for _, ch := range chunks {
		wg.Add(1)
		go func(ch chunk) {
			defer wg.Done()
			resultChan <- c.processChunk(data, ch)
		}(ch)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return c.mergeResults(resultChan)
}

func (c Uint32MmapParallel) processChunk(data []byte, ch chunk) IPv4set {
	seen := make(IPv4set, maxCapacity(ch.size()/avgIPv4size))
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

func (c Uint32MmapParallel) mergeResults(resultChan <-chan IPv4set) IPv4set {
	total := make(IPv4set)
	for result := range resultChan {
		for ip := range result {
			total[ip] = struct{}{}
		}
	}
	return total
}
