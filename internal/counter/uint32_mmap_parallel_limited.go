package counter

import (
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/dustin/go-humanize"
)

type Uint32MmapParallelLimited struct {
	Workers   int
	ChunkSize int
}

func (c Uint32MmapParallelLimited) Name() string {
	return fmt.Sprintf("uint32_mmap_parallel_limited(w=%d, cs=%s)", c.Workers, humanize.IBytes(uint64(c.ChunkSize)))
}

func (c Uint32MmapParallelLimited) Count(f *os.File) (int, error) {
	if c.Workers < 1 {
		return 0, fmt.Errorf("workers must be >= 1")
	}

	if c.ChunkSize < 1 {
		return 0, fmt.Errorf("chunk size must be >= 1")
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

func (c Uint32MmapParallelLimited) splitChunks(data []byte) []chunk {
	size := len(data)
	chunks := make([]chunk, 0, (size+c.ChunkSize-1)/c.ChunkSize)

	start := 0
	for start < size {
		end := start + c.ChunkSize
		if end > size {
			end = size
		}

		if end < size {
			for end < size && data[end] != '\n' || end < size {
				end++
			}
		}

		chunks = append(chunks, chunk{start: start, end: end})
		start = end
	}

	return chunks
}

func (c Uint32MmapParallelLimited) processChunksParallel(data []byte, chunks []chunk) IPv4set {
	chunkChan := make(chan int, len(chunks))
	for i := range chunks {
		chunkChan <- i
	}
	close(chunkChan)

	numWorkers := c.Workers
	if numWorkers > len(chunks) {
		numWorkers = len(chunks)
	}

	resultChan := make(chan IPv4set, numWorkers)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range chunkChan {
				resultChan <- c.processChunk(data, chunks[idx])
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return c.mergeResults(resultChan)
}

func (c Uint32MmapParallelLimited) processChunk(data []byte, ch chunk) IPv4set {
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

func (c Uint32MmapParallelLimited) mergeResults(resultChan <-chan IPv4set) IPv4set {
	total := make(IPv4set)
	for result := range resultChan {
		for ip := range result {
			total[ip] = struct{}{}
		}
	}
	return total
}
