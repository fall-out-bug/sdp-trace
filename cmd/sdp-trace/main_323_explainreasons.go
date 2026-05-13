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
