package main

import (
	"fmt"
	"io"
)

func explainReasons(reasons []string, stdout io.Writer) {
	for _, reason := range reasons {
		// Reasons are emitted verbatim as traceable verdict support, not as
		// prose-only interpretation.
		fmt.Fprintf(stdout, "Reason: %s\n", reason)
	}
}

func explainNextActions(actions []string, stdout io.Writer) {
	for _, action := range actions {
		// Next actions are remediation hints and do not upgrade the current
		// assessed state.
		fmt.Fprintf(stdout, "Next action: %s\n", action)
	}
}
