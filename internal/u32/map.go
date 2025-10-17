package u32

func NewMapSet(capacity int) Set {
	return &mapSet{
		m: make(map[uint32]struct{}, capacity),
	}
}

type mapSet struct {
	m map[uint32]struct{}
}

func (s *mapSet) Add(ip uint32) {
	s.m[ip] = struct{}{}
}

func (s *mapSet) Count() int {
	return len(s.m)
}
