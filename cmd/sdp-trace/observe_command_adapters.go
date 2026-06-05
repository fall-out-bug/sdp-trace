package main

import (
	"context"
	"io"
)

func runObserveCommand(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runObserve(args, stdout, stderr)
}

func runHarnessCommand(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runHarness(args, stdout, stderr)
}
