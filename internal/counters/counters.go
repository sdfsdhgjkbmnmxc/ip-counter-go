package counters

import (
	"os"

	"github.com/sdfsdhgjkbmnmxc/ip-counter-go/internal/u32"
)

type IPAddrCounter interface {
	Name() string
	Count(*os.File) (int, error)
}

var Registry = countersRegistry{
	NewMMapCounter("Map", func(size int64) u32.Set {
		return u32.NewMapSet(maxCapacity(int(size) / avgIPv4size))
	}),
	NewMMapCounter("Bitmap", func(int64) u32.Set {
		return u32.NewBitmapSet()
	}),
	NewParallelMMapCounter("ParallelBitmap", func(int64) u32.Set {
		return u32.NewAtomicBitmapSet()
	}),
	AutoCounter{
		SmallFiles: "Map",
		LargeFiles: "ParallelBitmap",
	},
	NaiveCounter{},
}

type countersRegistry []IPAddrCounter

func (l countersRegistry) Get(name string) IPAddrCounter {
	for _, counter := range l {
		if counter.Name() == name {
			return counter
		}
	}
	return nil
}

func (l countersRegistry) Help() string {
	help := "Available counters:\n"
	for _, counter := range l {
		help += " - " + counter.Name() + "\n"
	}
	return help
}
