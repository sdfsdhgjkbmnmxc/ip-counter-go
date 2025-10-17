package counter

import (
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/dustin/go-humanize"
	"golang.org/x/sync/errgroup"
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

	total, err := c.processChunksParallel(data, c.splitChunks(data))
	if err != nil {
		return 0, err
	}
	return total.Count(), nil
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

func (c Uint32MmapParallelLimited) processChunksParallel(data []byte, chunks []chunk) (uint32set, error) {
	chunkChan := make(chan int, len(chunks))
	for i := range chunks {
		chunkChan <- i
	}
	close(chunkChan)

	numWorkers := c.Workers
	if numWorkers > len(chunks) {
		numWorkers = len(chunks)
	}

	total := newIPv4Map(0)

	var mux sync.Mutex
	var g errgroup.Group

	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			for idx := range chunkChan {
				result, err := c.processChunk(data, chunks[idx])
				if err != nil {
					return err
				}

				mux.Lock()
				total = total.Union(result)
				mux.Unlock()
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return total, nil
}

func (c Uint32MmapParallelLimited) processChunk(data []byte, ch chunk) (uint32set, error) {
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
