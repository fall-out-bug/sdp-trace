package main

import (
	"fmt"
	"io"
)

func knownCheckpointCommand(command string, rest []string, stderr io.Writer) (string, []string, bool) {
	if command != "create" && command != "verify" {
		// Unknown checkpoint verbs are usage errors, not verifier states.
		fmt.Fprintln(stderr, "checkpoint requires create or verify")
		return "", nil, false
	}
	return command, rest, true
}
