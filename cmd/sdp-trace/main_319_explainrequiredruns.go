package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func explainRequiredRuns(requiredRuns []demo.RequiredRunResult, stdout io.Writer) {
	for _, requiredRun := range requiredRuns {
		// One stable line per required run keeps the human explanation auditable
		// without inventing a separate summary verdict.
		fmt.Fprintf(stdout, "Required run %s: %s\n", requiredRun.ID, requiredRun.State)
	}
}
