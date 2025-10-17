package counter

import (
	"bufio"
	"fmt"
	"os"
)

type BitmapCounter struct{}

func (c BitmapCounter) Name() string {
	return "bitmap"
}

func (c BitmapCounter) Count(f *os.File) (int, error) {
	bitmap := make([]byte, 1<<29)
	scanner := bufio.NewScanner(f)
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
