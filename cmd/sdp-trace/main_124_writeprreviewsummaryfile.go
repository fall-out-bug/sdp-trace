package main

import (
	"fmt"
	"io"
	"os"
)

func writePRReviewSummaryFile(path, summary string, stderr io.Writer) (int, bool) {
	if err := refuseExistingFile(path); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage, false
	}
	// Summaries are write-once CLI artifacts; refusing overwrite protects the
	// cited review text from accidental replacement.
	if err := os.WriteFile(path, []byte(summary), 0o644); err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	return 0, true
}
