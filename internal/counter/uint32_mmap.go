package counter

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

type Uint32Mmap struct{}

func (c Uint32Mmap) Name() string {
	return "uint32_mmap"
}

func (c Uint32Mmap) Count(r io.Reader) (int, error) {
	file, ok := r.(*os.File)
	if !ok {
		return 0, fmt.Errorf("mmap requires *os.File")
	}

	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	size := int(stat.Size())
	if size == 0 {
		return 0, nil
	}

	data, err := syscall.Mmap(int(file.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return 0, err
	}
	defer syscall.Munmap(data)

	seen := make(map[uint32]struct{})
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			if i > start {
				ip, err := parseIPv4FromBytes(data[start:i])
				if err != nil {
					return 0, fmt.Errorf("invalid IP address: %v", err)
				}
				seen[ip] = struct{}{}
			}
			start = i + 1
		}
	}

	if start < len(data) {
		ip, err := parseIPv4FromBytes(data[start:])
		if err != nil {
			return 0, fmt.Errorf("invalid IP address: %v", err)
		}
		seen[ip] = struct{}{}
	}

	return len(seen), nil
}
