package main

import (
	"context"
	"fmt"
	"io"
	"strings"
)

func runCommandSurface(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		fmt.Fprintf(stderr, "command-surface does not accept arguments: %s\n", strings.Join(args, " "))
		return exitUsage
	}
	if err := writeCommandSurfaceJSON(stdout); err != nil {
		return 2
	}
	return 0
}
