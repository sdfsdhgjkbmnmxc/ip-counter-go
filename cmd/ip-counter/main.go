package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sdfsdhgjkbmnmxc/ip-counter-go/internal/counter"
)

func main() {
	method := flag.String("method", "naive", "counting method")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-method=<method>] <filename>\n", os.Args[0])
		os.Exit(1)
	}

	c, ok := counter.Counters[*method]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown method: %s\n", *method)
		os.Exit(1)
	}

	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	count, err := c.Count(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting IPs: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(count)
}
