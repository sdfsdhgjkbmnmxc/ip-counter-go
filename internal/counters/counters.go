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
	NewMMapCounter("MapSet", func(fileSize int) u32.Set {
		return u32.NewMapSet(maxCapacity(fileSize / avgIPv4size))
	}),
	NewMMapCounter("BitmapSet", func(fileSize int) u32.Set {
		return u32.NewBitmapSet()
	}),
	NewMMapCounter("ComboSet", func(fileSize int) u32.Set {
		if fileSize/avgIPv4size > 28_000_000 { // see u32_test.go:TestMemoryCrossover
			return u32.NewBitmapSet()
		}
		return u32.NewMapSet(maxCapacity(fileSize / avgIPv4size))
	}),
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
