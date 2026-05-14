package main

import (
	"context"
	"io"
)

func runWrap(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	// Legacy wrap keeps parsing separate from recorder execution so malformed
	// wrapper metadata cannot create partial run artifacts.
	opts, command, code, ok := parseWrapArgs(args, stderr)
	if !ok {
		return code
	}
	return runLegacyWrapRecorder(ctx, opts, command, stdout, stderr)
}
