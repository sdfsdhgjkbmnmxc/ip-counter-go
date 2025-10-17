package u32

func BitmapSet() Set {
	return &bitmapSet{
		bitmap: make([]byte, 1<<29),
	}
}

type bitmapSet struct {
	count  int
	bitmap []byte
}

func (s *bitmapSet) Add(ip uint32) {
	byteIndex := ip >> 3
	bitIndex := ip & 7
	mask := byte(1 << bitIndex)

	if s.bitmap[byteIndex]&mask == 0 {
		s.count++
		s.bitmap[byteIndex] |= mask
	}
}

func (s *bitmapSet) Count() int {
	return s.count
}
