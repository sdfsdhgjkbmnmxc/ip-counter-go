package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sdfsdhgjkbmnmxc/ip-counter-go/internal/counters"
)

func main() {
	method := flag.String("method", "ComboSet", "counting method")
	flag.Parse()

	if flag.NArg() < 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [-method=<method>] <filename>\n", os.Args[0])
		os.Exit(1)
	}

	impl := counters.Registry.Get(*method)
	if impl == nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unknown method: %s. %s\n", *method, counters.Registry.Help())
		os.Exit(1)
	}

	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = file.Close() }()

	count, err := impl.Count(file)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error counting IPs: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(count)
}
