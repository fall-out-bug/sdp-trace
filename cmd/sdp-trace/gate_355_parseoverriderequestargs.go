package main

import (
	"fmt"
	"io"
)

func parseOverrideRequestArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	if !isOverrideRequest(args) {
		// The override namespace currently has one explicit write action.
		fmt.Fprintln(stderr, "override requires request")
		return nil, exitUsage, false
	}
	return parseOverrideRequestFlags(args[1:], stderr)
}
