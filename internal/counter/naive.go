package counter

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
)

type NaiveCounter struct {
	InitialBufferSize int
	MaxBufferSize     int
	Capacity          int
}

func (c NaiveCounter) Name() string {
	if c.InitialBufferSize == 0 && c.MaxBufferSize == 0 {
		return "naive"
	}
	return fmt.Sprintf("naive (init_buf=%s, max_buf=%s, cap=%d)",
		humanize.IBytes(uint64(c.InitialBufferSize)),
		humanize.IBytes(uint64(c.MaxBufferSize)),
		c.Capacity)
}

func (c NaiveCounter) Count(f *os.File) (int, error) {
	seen := make(map[string]struct{}, c.Capacity)
	scanner := bufio.NewScanner(f)

	if c.InitialBufferSize > 0 || c.MaxBufferSize > 0 {
		buf := make([]byte, 0, c.InitialBufferSize)
		scanner.Buffer(buf, c.MaxBufferSize)
	}

	for scanner.Scan() {
		ip := scanner.Text()
		seen[ip] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return len(seen), nil
}
