package counter

import (
	"fmt"
	"os"
	"syscall"

	"github.com/sdfsdhgjkbmnmxc/ip-counter-go/internal/u32"
)

type BitmapCounter struct{}

func (c BitmapCounter) Name() string { return "bitmap_counter" }

func (c BitmapCounter) Count(f *os.File) (int, error) {
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

	start := 0
	seen := u32.NewBitmapSet()

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
