package main

import (
	"fmt"
	"io"
	"strings"
)

func readRequiredProtectedInput(flag, path string, value any, stderr io.Writer) (int, bool) {
	if strings.TrimSpace(path) == "" {
		// Protected mode has no implicit local defaults for external trust
		// inputs.
		fmt.Fprintf(stderr, "protected gate requires %s\n", flag)
		return exitUsage, false
	}
	// All protected inputs are decoded as JSON artifacts before evaluation so
	// the gate never accepts unchecked path strings as trust evidence.
	if err := readJSONFile(path, value); err != nil {
		// Malformed trust inputs are usage/setup failures, not a green local gate
		// with omitted protected evidence.
		fmt.Fprintln(stderr, err)
		return exitUsage, false
	}
	return 0, true
}
