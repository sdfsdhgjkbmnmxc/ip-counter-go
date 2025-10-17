package u32

import (
	"runtime"
	"testing"

	"github.com/dustin/go-humanize"
)

func TestMemoryCrossover(t *testing.T) {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	beforeBitmap := m.HeapAlloc

	bitmap := NewBitmapSet()
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

		if used >= bitmapSize {
			t.Logf("CROSSOVER at %d elements: map %s >= bitmap %s", count, humanize.Bytes(used), humanize.Bytes(bitmapSize))
			break
		}

		if count >= 100_000_000 {
			t.Fatalf("CROSSOVER at %d elements: too many elements", count)
		}
	}
}
