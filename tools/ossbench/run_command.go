package main

import (
	"context"
	"os/exec"
)

func benchmarkCommand(ctx context.Context, cmd string, args []string, dir string) *exec.Cmd {
	c := exec.CommandContext(ctx, cmd, args...)
	if dir != "" {
		c.Dir = dir
	}
	return c
}
