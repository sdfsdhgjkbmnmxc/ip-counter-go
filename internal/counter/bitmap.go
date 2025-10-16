package counter

import (
	"bufio"
	"fmt"
	"io"
	"math/bits"
)

type BitmapCounter struct{}

func (c BitmapCounter) Name() string {
	return "bitmap"
}

func (c BitmapCounter) Count(r io.Reader) (int, error) {
	bitmap := make([]byte, 1<<29)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		ip, err := parseIPv4(scanner.Text())
		if err != nil {
			return 0, fmt.Errorf("invalid IP address: %v", err)
		}

		byteIndex := ip >> 3
		bitIndex := ip & 7
		bitmap[byteIndex] |= 1 << bitIndex
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	count := 0
	for _, b := range bitmap {
		count += bits.OnesCount8(b)
	}

	return count, nil
}
