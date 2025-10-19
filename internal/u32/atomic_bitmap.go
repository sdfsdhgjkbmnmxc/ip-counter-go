package u32

import "sync/atomic"

func NewAtomicBitmapSet() Set {
	return &atomicBitmapSet{
		bitmap: make([]atomic.Uint64, 1<<26), // 512MB / 8 bytes
	}
}

type atomicBitmapSet struct {
	count  atomic.Uint64
	bitmap []atomic.Uint64
}

func (s *atomicBitmapSet) Add(ip uint32) {
	index := ip >> 6    // divide by 64
	bitIndex := ip & 63 // modulo 64
	mask := uint64(1 << bitIndex)

	if s.bitmap[index].Load()&mask != 0 {
		return
	}

	for {
		old := s.bitmap[index].Load()
		if old&mask != 0 {
			return
		}
		if s.bitmap[index].CompareAndSwap(old, old|mask) {
			s.count.Add(1)
			return
		}
	}
}

func (s *atomicBitmapSet) Count() int {
	return int(s.count.Load())
}
