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
		{"roaring6", func() uint32set { return newIPv4Roaring(6) }},
		{"roaring8", func() uint32set { return newIPv4Roaring(8) }},
		{"roaring12", func() uint32set { return newIPv4Roaring(12) }},
		{"roaring16", func() uint32set { return newIPv4Roaring(16) }},
		{"roaring20", func() uint32set { return newIPv4Roaring(20) }},
		{"roaring24", func() uint32set { return newIPv4Roaring(24) }},
		{"bitmap", newIPv4Bitmap},
	}
}

func getSizes() []int {
	return []int{
		100,
		1_000,
		10_000,
		100_000,
		1_000_000,
		10_000_000,
	}
}

func BenchmarkAdd(b *testing.B) {
	for _, size := range getSizes() {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			for _, impl := range getImpl() {
				b.Run(impl.name, func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						s := impl.new()
						for j := 0; j < size; j++ {
							s.Add(uint32(j))
						}
						if s.Count() != size {
							b.Fatalf("Count() = %d, want %d", s.Count(), size)
						}
					}
				})
			}
		})
	}
}

func BenchmarkUnion(b *testing.B) {
	for _, size := range getSizes() {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			for _, impl := range getImpl() {
				b.Run(impl.name, func(b *testing.B) {
					b.ReportAllocs()

					s1 := impl.new()
					s2 := impl.new()

					b.StopTimer()
					for i := 0; i < size; i++ {
						s1.Add(uint32(i))
						s2.Add(uint32(i + size/2))
					}
					b.StartTimer()

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
