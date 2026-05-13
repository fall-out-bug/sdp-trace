package main

import (
	"io"
)

func runGateSubcommand(args []string, stdout, stderr io.Writer) (int, bool) {
	return runOptionalSubcommand(args, stdout, stderr, gateSubcommandHandlers)
}
