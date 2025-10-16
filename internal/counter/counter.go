package counter

import "io"

type IPAddrCounter interface {
	Count(io.Reader) (int64, error)
}

var Counters = map[string]IPAddrCounter{
	"naive": &NaiveCounter{
		InitialBufferSize: 64 * 1024,
		MaxBufferSize:     1024 * 1024,
	},
}
