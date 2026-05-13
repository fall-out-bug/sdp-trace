package main

import (
	"context"
	"io"
)

func runDryRun(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runPreviewCommand("dry-run", "simulation", args, stdout, stderr)
}
