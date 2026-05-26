package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func runWrapCommand(ctx context.Context, bin string, args []string, dir string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	stdout, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("wrap failed: %v\nstderr: %s", err, strings.TrimSpace(string(wrapStderr(err))))
	}
	return stdout, nil
}
