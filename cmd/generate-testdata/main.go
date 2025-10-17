package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s <count>\n", os.Args[0])
		os.Exit(1)
	}

	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Invalid count: %v\n", err)
		os.Exit(1)
	}

	for i := 0; i < count; i++ {
		fmt.Printf("%d.%d.%d.%d\n",
			rand.IntN(256),
			rand.IntN(256),
			rand.IntN(256),
			rand.IntN(256),
		)
	}
}
