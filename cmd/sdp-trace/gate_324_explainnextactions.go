package main

import (
	"fmt"
	"io"
)

func explainNextActions(actions []string, stdout io.Writer) {
	for _, action := range actions {
		// Next actions are remediation hints and do not upgrade the current
		// assessed state.
		fmt.Fprintf(stdout, "Next action: %s\n", action)
	}
}
