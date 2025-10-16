package counter

import "io"

type IPAddrCounter interface {
	Name() string
	Count(io.Reader) (int64, error)
}

var Counters = []IPAddrCounter{
	&NaiveCounter{InitialBufferSize: 64 * 1024, MaxBufferSize: 4096}, // go defaults
	&NaiveCounter{InitialBufferSize: 64 * 1024, MaxBufferSize: 16},
	&NaiveCounter{InitialBufferSize: 128 * 1024, MaxBufferSize: 16},
}
