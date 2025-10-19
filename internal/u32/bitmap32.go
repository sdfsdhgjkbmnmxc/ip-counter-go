package u32

func NewBitmap32Set() Set {
	return &bitmap32Set{
		bitmap: make([]uint32, 1<<27), // 512MB / 4 bytes
	}
}

type bitmap32Set struct {
	count  int
	bitmap []uint32
}

func (s *bitmap32Set) Add(ip uint32) {
	index := ip >> 5    // divide by 32
	bitIndex := ip & 31 // modulo 32
	mask := uint32(1 << bitIndex)

	if s.bitmap[index]&mask == 0 {
		s.count++
		s.bitmap[index] |= mask
	}
}

func (s *bitmap32Set) Count() int {
	return s.count
}
