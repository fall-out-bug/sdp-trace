package prreview

import (
	"errors"
	"os/exec"
	"strings"
)

func runnerUnavailable(err error) bool {
	return errors.Is(err, exec.ErrNotFound) || runnerNotFoundMessage(err)
}

func runnerNotFoundMessage(err error) bool {
	return strings.Contains(err.Error(), "executable file not found")
}
