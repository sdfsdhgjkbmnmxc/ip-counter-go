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
	seen := newIPv4Map(0)
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
