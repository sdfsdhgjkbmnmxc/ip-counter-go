package counter

import (
	"bufio"
	"fmt"
	"os"
)

type Uint32Map struct{}

func (c Uint32Map) Name() string {
	return "uint32_map"
}

func (c Uint32Map) Count(f *os.File) (int, error) {
	seen := make(map[uint32]struct{})
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		ip, err := parseIPv4(scanner.Text())
		if err != nil {
			return 0, fmt.Errorf("invalid IP address: %v", err)
		}
		seen[ip] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return len(seen), nil
}
