package main

import (
	"fmt"
	"io"
)

func checkpointCommand(args []string, stderr io.Writer) (string, []string, bool) {
	if len(args) == 0 {
		// A checkpoint command without a verb cannot decide whether it is
		// creating signing evidence or replaying it.
		fmt.Fprintln(stderr, "checkpoint requires create or verify")
		return "", nil, false
	}
	// Keep the remainder untouched so each subcommand owns its flag contract.
	return knownCheckpointCommand(args[0], args[1:], stderr)
}
