package u32

func NewBitmapSet() Set {
	return &bitmapSet{
		bitmap: make([]byte, 1<<29), // 512MB / 1 byte
	}
}

type bitmapSet struct {
	bitmap []byte
}

func (s *bitmapSet) Add(ip uint32) bool {
	byteIndex := ip >> 3 // divide by 8
	bitIndex := ip & 7   // modulo 8
	mask := byte(1 << bitIndex)

	if s.bitmap[byteIndex]&mask == 0 {
		s.bitmap[byteIndex] |= mask
		return true
	}

	return false
}
