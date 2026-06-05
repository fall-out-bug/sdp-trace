package main

import (
	"io"
	"strings"
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

func isHelpArg(arg string) bool {
	return arg == "help" || arg == "--help" || arg == "-h"
}

func subcommandName(label string) string {
	if before, _, ok := strings.Cut(label, " "); ok {
		// Help labels can include usage suffixes; dispatch diagnostics should
		// name only the stable subcommand token.
		return before
	}
	return label
}
