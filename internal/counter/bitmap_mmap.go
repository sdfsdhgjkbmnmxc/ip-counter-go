package counter

import (
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

	return c.processData(data), nil
}

func (c BitmapMmap) processData(data []byte) int {
	bitmap := make([]byte, 1<<29)
	count := 0
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			if i > start {
				if c.setBit(bitmap, data[start:i]) {
					count++
				}
			}
			start = i + 1
		}
	}

	if start < len(data) {
		if c.setBit(bitmap, data[start:]) {
			count++
		}
	}

	return count
}

func (c BitmapMmap) setBit(bitmap []byte, line []byte) bool {
	ip, err := parseIPv4FromBytes(line)
	if err != nil {
		return false
	}

	byteIndex := ip >> 3
	bitIndex := ip & 7
	mask := byte(1 << bitIndex)

	if bitmap[byteIndex]&mask == 0 {
		bitmap[byteIndex] |= mask
		return true
	}

	return false
}
