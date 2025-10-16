package counter

import (
	"bufio"
	"fmt"
	"io"
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
	return fmt.Sprintf("naive (init buf: %d, max buf: %d, capacity: %d)", c.InitialBufferSize, c.MaxBufferSize, c.Capacity)
}

func (c NaiveCounter) Count(r io.Reader) (int, error) {
	seen := make(map[string]struct{}, c.Capacity)
	scanner := bufio.NewScanner(r)

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
