package counters

import (
	"errors"
	"fmt"
)

var ErrInvalidIP = errors.New("invalid IP address")

func wrapInvalidIPError(err error) error {
	return fmt.Errorf("%w: %v", ErrInvalidIP, err)
}
