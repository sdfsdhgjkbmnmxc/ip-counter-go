package counter

import (
	"bufio"
	"io"
)

type NaiveCounter struct {
	InitialBufferSize int
	MaxBufferSize     int
}

func (c *NaiveCounter) Count(r io.Reader) (int64, error) {
	seen := make(map[string]struct{})
	scanner := bufio.NewScanner(r)

	buf := make([]byte, 0, c.InitialBufferSize)
	scanner.Buffer(buf, c.MaxBufferSize)

	for scanner.Scan() {
		ip := scanner.Text()
		seen[ip] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return int64(len(seen)), nil
}
