package main

import (
	"context"
	"io"
)

func runCommandSurface(_ context.Context, _ []string, stdout, _ io.Writer) int {
	if err := writeCommandSurfaceJSON(stdout); err != nil {
		return 2
	}
	return 0
}
