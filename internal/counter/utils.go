package counter

import "fmt"

const (
	avgIPv4size   = 15
	maxIPv4number = 1 << 32
)

func maxCapacity(v int) int {
	if v > maxIPv4number {
		return maxIPv4number
	}
	return v
}

func parseIPv4(s string) (uint32, error) {
	return parseIPv4FromBytes([]byte(s))
}

func parseIPv4FromBytes(s []byte) (uint32, error) {
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

type chunk struct {
	start, end int
}

func (c chunk) size() int {
	return c.end - c.start
}
