package counter

import (
	"fmt"
	"os"
	"syscall"
)

type BitmapMmap struct{}

func (c BitmapMmap) Name() string {
	return "bitmap_mmap"
}

func (c BitmapMmap) Count(f *os.File) (int, error) {
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}

	size := int(stat.Size())
	if size == 0 {
		return 0, nil
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return 0, err
	}
	defer func() { _ = syscall.Munmap(data) }()

	seen := newIPv4BitSet()
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			if i > start {
				ip, err := parseIPv4FromBytes(data[start:i])
				if err != nil {
					return 0, fmt.Errorf("invalid IP address: %v", err)
				}
				seen.Add(ip)
			}
			start = i + 1
		}
	}

	if start < len(data) {
		ip, err := parseIPv4FromBytes(data[start:])
		if err != nil {
			return 0, fmt.Errorf("invalid IP address: %v", err)
		}
		seen.Add(ip)
	}

	return seen.Count(), nil
}
