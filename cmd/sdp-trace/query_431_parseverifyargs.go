package main

import (
	"fmt"
	"io"
)

func parseVerifyArgs(args []string, stderr io.Writer) (string, int, bool) {
	if len(args) == 0 {
		// Verify requires a concrete retained run directory.
		fmt.Fprintln(stderr, "verify requires <run-dir>")
		return "", exitUsage, false
	}
	runDir := args[0]
	if !existingDirectory(runDir) {
		// Missing run roots are cannot_verify rather than usage once a path was
		// supplied.
		fmt.Fprintf(stderr, "run directory does not exist: %s\n", runDir)
		return "", exitCannotVerify, false
	}
	return runDir, 0, true
}
