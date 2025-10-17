package counter

import (
	"fmt"
	"testing"
)

type impl struct {
	name string
	new  func() uint32set
}

func getImpl() []impl {
	return []impl{
		{"map", func() uint32set { return newIPv4Map(0) }},
		{"bitset", newIPv4BitSet},
	}
}

func TestUint32SetBasics(t *testing.T) {
	for _, impl := range getImpl() {
		t.Run(impl.name, func(t *testing.T) {
			s := impl.new()

			s.Add(0x01020304)
			if s.Count() != 1 {
				t.Errorf("Count() = %d, want 1", s.Count())
			}

			s.Add(0x01020304)
			if s.Count() != 1 {
				t.Errorf("Count() after duplicate = %d, want 1", s.Count())
			}

			s.Add(0x05060708)
			if s.Count() != 2 {
				t.Errorf("Count() = %d, want 2", s.Count())
			}
		})
	}
}

func TestUint32SetUnion(t *testing.T) {
	for _, impl := range getImpl() {
		t.Run(impl.name, func(t *testing.T) {
			s1 := impl.new()
			s1.Add(0x01020304)
			s1.Add(0x05060708)

			s2 := impl.new()
			s2.Add(0x05060708)
			s2.Add(0x090A0B0C)

			result := s1.Union(s2)
			if result.Count() != 3 {
				t.Errorf("Union Count() = %d, want 3", result.Count())
			}
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	for _, impl := range getImpl() {
		b.Run(impl.name, func(b *testing.B) {
			s := impl.new()
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				s.Add(uint32(i))
			}
		})
	}
}

func BenchmarkUnion(b *testing.B) {
	for _, size := range []int{
		100,
		1000,
		10000,
		100000,
	} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			for _, impl := range getImpl() {
				b.Run(impl.name, func(b *testing.B) {
					s1 := impl.new()
					s2 := impl.new()

					for i := 0; i < size; i++ {
						s1.Add(uint32(i))
						s2.Add(uint32(i + size/2))
					}

					b.ReportAllocs()
					b.ResetTimer()

					for i := 0; i < b.N; i++ {
						s1Copy := impl.new()
						s1Copy.Union(s1)
						s1Copy.Union(s2)
					}
				})
			}
		})
	}
}
