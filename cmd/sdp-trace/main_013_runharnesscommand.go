package main

import (
	"context"
	"io"
)

func runHarnessCommand(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runHarness(args, stdout, stderr)
}
