package counter

import (
	"bufio"
	"fmt"
	"io"
)

type Uint32Map struct{}

func (c *Uint32Map) Name() string {
	return "uint32_map"
}

func (c *Uint32Map) Count(r io.Reader) (int, error) {
	seen := make(map[uint32]struct{})
	scanner := bufio.NewScanner(r)

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

func parseIPv4(s string) (uint32, error) {
	var ip uint32
	var octet uint32
	var shift = 24
	var dotCount int
	var hasDigit bool

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			octet = octet*10 + uint32(c-'0')
			if octet > 255 {
				return 0, fmt.Errorf("invalid octet")
			}
			hasDigit = true
		} else if c == '.' {
			if !hasDigit {
				return 0, fmt.Errorf("empty octet")
			}
			if shift < 0 {
				return 0, fmt.Errorf("too many dots")
			}
			ip |= octet << shift
			shift -= 8
			octet = 0
			dotCount++
			hasDigit = false
		} else {
			return 0, fmt.Errorf("invalid character")
		}
	}

	if !hasDigit {
		return 0, fmt.Errorf("empty octet")
	}

	if dotCount != 3 {
		return 0, fmt.Errorf("invalid IP format")
	}

	ip |= octet

	return ip, nil
}
