package counter

import (
	"bufio"
	"fmt"
	"io"
	"net/netip"
)

type NetIPMap struct{}

func (c *NetIPMap) Name() string {
	return "netip_map"
}

func (c *NetIPMap) Count(r io.Reader) (int, error) {
	seen := make(map[netip.Addr]struct{})
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		ip, err := netip.ParseAddr(scanner.Text())
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
