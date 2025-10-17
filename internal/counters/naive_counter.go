package counters

import (
	"bufio"
	"os"
)

type NaiveCounter struct{}

func (c NaiveCounter) Name() string { return "NaiveCounter" }

func (c NaiveCounter) Count(f *os.File) (int, error) {
	seen := make(map[string]struct{})
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		ip := scanner.Text()
		seen[ip] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return len(seen), nil
}
