package u32

func NewMapSet(capacity int) Set {
	return &mapSet{
		m: make(map[uint32]struct{}, capacity),
	}
}

type mapSet struct {
	m map[uint32]struct{}
}

func (s *mapSet) Add(ip uint32) bool {
	if _, exists := s.m[ip]; exists {
		return false
	}
	s.m[ip] = struct{}{}
	return true
}
