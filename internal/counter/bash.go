package counter

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type BashSort struct{}

func (c BashSort) Name() string {
	return "bash_sort"
}

func (c BashSort) Count(r io.Reader) (int, error) {
	file, ok := r.(*os.File)
	if !ok {
		return 0, fmt.Errorf("bash sort requires *os.File")
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("sort -u %s | wc -l", file.Name()))
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
