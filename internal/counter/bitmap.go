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
	seen := newIPv4Bitmap()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		ip, err := parseIPv4(scanner.Text())
		if err != nil {
			return 0, fmt.Errorf("invalid IP address: %v", err)
		}
		seen.Add(ip)
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return seen.Count(), nil
}
