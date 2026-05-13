package main

import (
	"io"
)

func runOptionalSubcommand(args []string, stdout, stderr io.Writer, handlers map[string]subcommandHandler) (int, bool) {
	if len(args) == 0 {
		return 0, false
	}
	handler, ok := handlers[args[0]]
	if !ok {
		// Optional dispatch lets parent commands fall back to flag parsing.
		return 0, false
	}
	return handler(args[1:], stdout, stderr), true
}
