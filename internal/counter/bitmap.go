package counter

import (
	"bufio"
	"fmt"
	"io"
)

type BitmapCounter struct{}

func (c BitmapCounter) Name() string {
	return "bitmap"
}

func (c BitmapCounter) Count(r io.Reader) (int, error) {
	bitmap := make([]byte, 1<<29)
	scanner := bufio.NewScanner(r)
	count := 0

	for scanner.Scan() {
		ip, err := parseIPv4(scanner.Text())
		if err != nil {
			return 0, fmt.Errorf("invalid IP address: %v", err)
		}

		byteIndex := ip >> 3
		bitIndex := ip & 7
		mask := byte(1 << bitIndex)

		if bitmap[byteIndex]&mask == 0 {
			count++
			bitmap[byteIndex] |= mask
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}
