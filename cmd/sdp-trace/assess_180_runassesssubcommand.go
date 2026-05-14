package main

import (
	"io"
)

func runAssessSubcommand(args []string, stdout, stderr io.Writer) (int, bool) {
	return runOptionalSubcommand(args, stdout, stderr, assessSubcommandHandlers)
}
