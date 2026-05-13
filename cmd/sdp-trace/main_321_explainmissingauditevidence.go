package main

import (
	"fmt"
	"io"
)

func explainMissingAuditEvidence(missingEvidence []string, stdout io.Writer) {
	for _, missing := range missingEvidence {
		// Missing audit evidence stays visible as a concrete gap; explanation
		// output must not collapse it into a green summary.
		fmt.Fprintf(stdout, "Missing audit evidence: %s\n", missing)
	}
}
