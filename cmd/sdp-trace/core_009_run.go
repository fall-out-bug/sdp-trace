package main

import (
	"io"
)

func run(args []string, stdout, stderr io.Writer) int {
	if topLevelHelp(args) {
		printUsage(stdout)
		return 0
	}
	// The first token is the only command selector; everything after it stays
	// command-owned so subcommands can preserve their own evidence contract.
	return dispatchCommand(args[0], args[1:], stdout, stderr)
}
