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

type ipv4bitset struct {
	bits  uint32
	count int
}

func (s ipv4bitset) Add(ip uint32) {
	mask := uint32(1) << (ip % 32)
	if s.bits&mask == 0 {
		s.bits |= mask
		s.count++
	}
}

func (s ipv4bitset) Union(other uint32set) uint32set {
	switch other.(type) {
	case ipv4bitset:
		otherSet := other.(ipv4bitset)
		newBits := s.bits | otherSet.bits
		newCount := 0
		for i := 0; i < 32; i++ {
			if newBits&(1<<i) != 0 {
				newCount++
			}
		}
		return ipv4bitset{bits: newBits, count: newCount}
	default:
		panic("unreachable")
	}
}

func (s ipv4bitset) Count() int {
	return s.count
}
