package counter

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type BashCounter struct{}

func (c BashCounter) Name() string {
	return "BashCounter"
}

func (c BashCounter) Count(f *os.File) (int, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("sort -u %s | wc -l", f.Name()))
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, err
	}

	return count, nil
}
