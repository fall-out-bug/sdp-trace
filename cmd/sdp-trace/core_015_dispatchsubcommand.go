package main

import (
	"fmt"
	"io"
)

func dispatchSubcommand(cmd string, args []string, stdout, stderr io.Writer, label string, handlers map[string]subcommandHandler) int {
	handler, ok := handlers[cmd]
	if !ok {
		// Keep command-family names stable in diagnostics even when the usage label
		// contains argument suffixes.
		fmt.Fprintf(stderr, "unknown %s command: %s\n", subcommandName(label), cmd)
		return exitUsage
	}
	return handler(args, stdout, stderr)
}
