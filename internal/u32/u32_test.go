package u32

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/dustin/go-humanize"
)

// BenchmarkAdd/size-100/MapSet-10  	 3258494	      1114 ns/op	    1248 B/op	       5 allocs/op
// BenchmarkAdd/size-100/BitmapSet-10         	     666	   5416677 ns/op	536870957 B/op	       2 allocs/op
// BenchmarkAdd/size-1000/MapSet-10           	  298838	     11303 ns/op	   19080 B/op	       7 allocs/op
// BenchmarkAdd/size-1000/BitmapSet-10        	     634	   5397548 ns/op	536870960 B/op	       2 allocs/op
// BenchmarkAdd/size-10000/MapSet-10          	   31196	    115934 ns/op	  152248 B/op	      35 allocs/op
// BenchmarkAdd/size-10000/BitmapSet-10       	     630	   5449612 ns/op	536870962 B/op	       2 allocs/op
// BenchmarkAdd/size-100000/MapSet-10         	    2198	   1480138 ns/op	 1217729 B/op	     259 allocs/op
// BenchmarkAdd/size-100000/BitmapSet-10      	     612	   5632465 ns/op	536870962 B/op	       2 allocs/op
// BenchmarkAdd/size-1000000/MapSet-10        	     103	  33824183 ns/op	19482710 B/op	    4099 allocs/op
// BenchmarkAdd/size-1000000/BitmapSet-10     	     444	   7792310 ns/op	536870963 B/op	       2 allocs/op
// BenchmarkAdd/size-2000000/MapSet-10        	      34	  96708476 ns/op	38961250 B/op	    8195 allocs/op
// BenchmarkAdd/size-2000000/BitmapSet-10     	     342	  10038622 ns/op	536870959 B/op	       2 allocs/op
// BenchmarkAdd/size-3000000/MapSet-10        	      22	 155018648 ns/op	38961238 B/op	    8195 allocs/op
// BenchmarkAdd/size-3000000/BitmapSet-10     	     283	  12478577 ns/op	536870959 B/op	       2 allocs/op
// BenchmarkAdd/size-4000000/MapSet-10        	      14	 226973720 ns/op	77922401 B/op	   16387 allocs/op
// BenchmarkAdd/size-4000000/BitmapSet-10     	     240	  14818751 ns/op	536870961 B/op	       2 allocs/op
// BenchmarkAdd/size-5000000/MapSet-10        	      12	 291973705 ns/op	77922424 B/op	   16387 allocs/op
// BenchmarkAdd/size-5000000/BitmapSet-10     	     207	  17102866 ns/op	536870954 B/op	       2 allocs/op
// BenchmarkAdd/size-6000000/MapSet-10        	       9	 356582931 ns/op	77922413 B/op	   16387 allocs/op
// BenchmarkAdd/size-6000000/BitmapSet-10     	     182	  19351185 ns/op	536870955 B/op	       2 allocs/op
// BenchmarkAdd/size-7000000/MapSet-10        	       7	 447814381 ns/op	89222030 B/op	   18738 allocs/op
// BenchmarkAdd/size-7000000/BitmapSet-10     	     164	  21777815 ns/op	536870965 B/op	       2 allocs/op
// BenchmarkAdd/size-8000000/MapSet-10        	       6	 508038562 ns/op	155844696 B/op	   32771 allocs/op
// BenchmarkAdd/size-8000000/BitmapSet-10     	     147	  24018921 ns/op	536870959 B/op	       2 allocs/op
// BenchmarkAdd/size-9000000/MapSet-10        	       6	 582489160 ns/op	155844696 B/op	   32771 allocs/op
// BenchmarkAdd/size-9000000/BitmapSet-10     	     134	  26479166 ns/op	536870963 B/op	       2 allocs/op
// BenchmarkAdd/size-10000000/MapSet-10       	       5	 655795292 ns/op	155844721 B/op	   32771 allocs/op
// BenchmarkAdd/size-10000000/BitmapSet-10    	     124	  28725113 ns/op	536870958 B/op	       2 allocs/op
func BenchmarkAdd(b *testing.B) {
	for _, size := range []int{
		100,
		1_000,
		10_000,
		100_000,
		1_000_000,
		2_000_000,
		3_000_000,
		4_000_000,
		5_000_000,
		6_000_000,
		7_000_000,
		8_000_000,
		9_000_000,
		10_000_000,
	} {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			for _, tc := range []struct {
				name string
				new  func() Set
			}{
				{
					name: "MapSet",
					new:  func() Set { return NewMapSet(size) },
				},
				{
					name: "BitmapSet",
					new:  func() Set { return BitmapSet() },
				},
			} {
				b.Run(tc.name, func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						s := tc.new()
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

func TestMemoryCrossover(t *testing.T) {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	beforeBitmap := m.HeapAlloc

	bitmap := BitmapSet()
	bitmap.Add(1)

	runtime.ReadMemStats(&m)
	bitmapSize := m.HeapAlloc - beforeBitmap
	t.Logf("Bitmap actual size: %s", humanize.Bytes(bitmapSize))

	runtime.GC()
	runtime.ReadMemStats(&m)
	before := m.HeapAlloc

	s := NewMapSet(0)
	count := 0

	for {
		for i := 0; i < 100_000; i++ {
			s.Add(uint32(count))
			count++
		}

		runtime.ReadMemStats(&m)
		used := m.HeapAlloc - before

		t.Logf("Elements: %d, Memory: %s (%.2f bytes/element)", count, humanize.Bytes(used), float64(used)/float64(count))

		if used >= bitmapSize {
			t.Logf("CROSSOVER at %d elements: map %s >= bitmap %s", count, humanize.Bytes(used), humanize.Bytes(bitmapSize))
			break
		}

		if count >= 100_000_000 {
			t.Fatalf("CROSSOVER at %d elements: too many elements", count)
		}
	}
}
