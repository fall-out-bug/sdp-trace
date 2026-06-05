package main

import (
	"context"
	"io"
)

func runWrappedCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts, command, code, ok := parseWrappedCommandArgs(args, stderr)
	if !ok {
		return code
	}
	// The modern run command requires an explicit task and contract choice
	// before recorder execution.
	return runTaskRecorder(ctx, opts, command, stdout, stderr)
}
