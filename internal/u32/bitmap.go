package u32

import "sync"

func NewBitmapSet() Set {
	return &bitmapSet{
		bitmap: make([]byte, 1<<29), // 512MB / 1 byte
	}
}

type bitmapSet struct {
	count  int
	bitmap []byte
}

func (s *bitmapSet) Add(ip uint32) {
	byteIndex := ip >> 3 // divide by 8
	bitIndex := ip & 7   // modulo 8
	mask := byte(1 << bitIndex)

	if s.bitmap[byteIndex]&mask == 0 {
		s.count++
		s.bitmap[byteIndex] |= mask
	}
}

func (s *bitmapSet) Count() int {
	return s.count
}

func NewSyncBitmapSet() Set {
	return &syncBitmapSet{
		bitmap: make([]byte, 1<<29), // 512MB / 1 byte
	}
}

type syncBitmapSet struct {
	count  int
	bitmap []byte
	mux    sync.Mutex
}

func (s *syncBitmapSet) Add(ip uint32) {
	byteIndex := ip >> 3 // divide by 8
	bitIndex := ip & 7   // modulo 8
	mask := byte(1 << bitIndex)

	s.mux.Lock()
	defer s.mux.Unlock()

	if s.bitmap[byteIndex]&mask == 0 {
		s.count++
		s.bitmap[byteIndex] |= mask
	}
}

func (s *syncBitmapSet) Count() int {
	return s.count
}
