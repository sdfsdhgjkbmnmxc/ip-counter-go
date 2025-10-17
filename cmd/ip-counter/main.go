package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sdfsdhgjkbmnmxc/ip-counter-go/internal/counter"
)

func main() {
	method := flag.String("method", counter.MapCounter{}.Name(), "counting method")
	flag.Parse()

	if flag.NArg() < 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [-method=<method>] <filename>\n", os.Args[0])
		os.Exit(1)
	}

	var methodImpl counter.IPAddrCounter
	for _, c := range counter.Counters {
		if c.Name() == *method {
			methodImpl = c
			break
		}
	}

	if methodImpl == nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unknown method: %s\n", *method)
		os.Exit(1)
	}

	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = file.Close() }()

	count, err := methodImpl.Count(file)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error counting IPs: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(count)
}
