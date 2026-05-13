package main

import (
	"context"
	"fmt"
	"io"
)

func dispatchCommand(cmd string, args []string, stdout, stderr io.Writer) int {
	handler, ok := commandHandlers[cmd]
	if !ok {
		// Unknown commands are usage defects, not verifier or gate verdicts.
		fmt.Fprintf(stderr, "unknown command: %s\n", cmd)
		printUsage(stderr)
		return 1
	}
	return handler(context.Background(), args, stdout, stderr)
}
