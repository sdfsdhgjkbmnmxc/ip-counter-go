package counter

import (
	"os"
)

type IPAddrCounter interface {
	Name() string
	Count(*os.File) (int, error)
}

var Counters = []IPAddrCounter{
	BashCounter{},
	BitmapCounter{},
	MapCounter{},
	NaiveCounter{},
}
