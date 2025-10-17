package counter

import "os"

type IPAddrCounter interface {
	Name() string
	Count(*os.File) (int, error)
}

var Counters = []IPAddrCounter{
	&NaiveCounter{InitialBufferSize: 64 * 1024, MaxBufferSize: 4096}, // go defaults
	&NaiveCounter{InitialBufferSize: 64 * 1024, MaxBufferSize: 16, Capacity: 1024},
	&NetIPMap{},
	&Uint32Map{},
	&Uint32Mmap{},
	&Uint32MmapParallel{Workers: 2},
	&Uint32MmapParallel{Workers: 4},
	&Uint32MmapParallel{Workers: 8},
	&Uint32MmapParallel{Workers: 16},
	&BitmapCounter{},
	&BitmapMmap{},
	&BashSort{},
}
