package counter

import "sync"

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
	m   map[uint32]struct{}
	mux sync.Mutex
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

func newIPv4BitSet() uint32set {
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
