package counter

import (
	"os"
	"runtime"

	"github.com/dustin/go-humanize"
)

type IPAddrCounter interface {
	Name() string
	Count(*os.File) (int, error)
}

var Counters = []IPAddrCounter{
	&NaiveCounter{InitialBufferSize: 64 * humanize.KiByte, MaxBufferSize: 4 * humanize.KiByte}, // go defaults
	&NaiveCounter{InitialBufferSize: 64 * humanize.KiByte, MaxBufferSize: 16, Capacity: 1024},
	&NetIPMap{},
	&Uint32Map{},
	&Uint32Mmap{},
	&Uint32MmapParallel{Workers: 2},
	&Uint32MmapParallel{Workers: runtime.NumCPU()},
	&Uint32MmapParallelLimited{Workers: runtime.NumCPU(), ChunkSize: 32 * humanize.KiByte},
	&Uint32MmapParallelLimited{Workers: runtime.NumCPU(), ChunkSize: 32 * humanize.MiByte},
	&BitmapCounter{},
	&BitmapMmap{},
	&BitmapMmapRoaring{},
	&BashSort{},
}
