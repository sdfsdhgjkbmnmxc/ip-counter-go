package counter

type uint32set interface {
	Add(ip uint32)
	Count() int
	Union(other uint32set) uint32set
}

func newIPv4Map(capacity int) uint32set {
	return &ipv4map{
		m: make(map[uint32]struct{}, capacity),
	}
}

type ipv4map struct {
	m map[uint32]struct{}
}

func (s *ipv4map) Add(ip uint32) {
	s.m[ip] = struct{}{}
}

func (s *ipv4map) Union(other uint32set) uint32set {
	switch other.(type) {
	case *ipv4map:
		for ip := range other.(*ipv4map).m {
			s.Add(ip)
		}
	default:
		panic("unreachable")
	}

	return s
}

func (s *ipv4map) Count() int {
	return len(s.m)
}

func newIPv4Bitmap() uint32set {
	return &ipv4bitmap{
		bitmap: make([]byte, 1<<29),
	}
}

type ipv4bitmap struct {
	count  int
	bitmap []byte
}

func (s *ipv4bitmap) Add(ip uint32) {
	byteIndex := ip >> 3
	bitIndex := ip & 7
	mask := byte(1 << bitIndex)

	if s.bitmap[byteIndex]&mask == 0 {
		s.count++
		s.bitmap[byteIndex] |= mask
	}
}

func (s *ipv4bitmap) Union(other uint32set) uint32set {
	switch other.(type) {
	case *ipv4bitmap:
		otherSet := other.(*ipv4bitmap)
		for i := 0; i < len(s.bitmap); i++ {
			combined := s.bitmap[i] | otherSet.bitmap[i]
			for b := byte(1); b != 0; b <<= 1 {
				if combined&b != 0 && s.bitmap[i]&b == 0 {
					s.count++
				}
			}
			s.bitmap[i] = combined
		}
	default:
		panic("unreachable")
	}
	return s
}

func (s *ipv4bitmap) Count() int {
	return s.count
}

func newIPv4Roaring(segmentBits int) uint32set {
	return &ipv4roaring{
		segments:    make(map[uint32][]byte),
		segmentBits: segmentBits,
		segmentSize: (1 << segmentBits) / 8,
	}
}

type ipv4roaring struct {
	segments    map[uint32][]byte
	segmentBits int
	segmentSize int
	count       int
}

func (s *ipv4roaring) Add(ip uint32) {
	high := uint32(ip >> s.segmentBits)
	low := uint32(ip & ((1 << s.segmentBits) - 1))

	seg := s.segments[high]
	if seg == nil {
		seg = make([]byte, s.segmentSize)
		s.segments[high] = seg
	}

	byteIndex := low >> 3
	bitIndex := low & 7

	if byteIndex >= uint32(s.segmentSize) {
		panic("byteIndex out of range")
	}

	mask := byte(1 << bitIndex)

	if seg[byteIndex]&mask == 0 {
		s.count++
		seg[byteIndex] |= mask
	}
}

func (s *ipv4roaring) Union(other uint32set) uint32set {
	switch other.(type) {
	case *ipv4roaring:
		otherSet := other.(*ipv4roaring)
		for high, otherChunk := range otherSet.segments {
			seg := s.segments[high]
			if seg == nil {
				seg = make([]byte, s.segmentSize)
				copy(seg, otherChunk)
				s.segments[high] = seg
				for _, b := range seg {
					s.count += popcount(b)
				}
			} else {
				for i := 0; i < len(seg); i++ {
					combined := seg[i] | otherChunk[i]
					if combined != seg[i] {
						s.count += popcount(combined) - popcount(seg[i])
					}
					seg[i] = combined
				}
			}
		}
	default:
		panic("unreachable")
	}
	return s
}

func (s *ipv4roaring) Count() int {
	return s.count
}

func popcount(b byte) int {
	count := 0
	for b != 0 {
		count++
		b &= b - 1
	}
	return count
}
