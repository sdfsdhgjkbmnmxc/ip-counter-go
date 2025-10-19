package u32

func NewBitmap64Set() Set {
	return &bitmap64Set{
		bitmap: make([]uint64, 1<<26), // 512MB / 8 bytes
	}
}

type bitmap64Set struct {
	count  int
	bitmap []uint64
}

func (s *bitmap64Set) Add(ip uint32) {
	index := ip >> 6    // divide by 64
	bitIndex := ip & 63 // modulo 64
	mask := uint64(1 << bitIndex)

	if s.bitmap[index]&mask == 0 {
		s.count++
		s.bitmap[index] |= mask
	}
}

func (s *bitmap64Set) Count() int {
	return s.count
}
