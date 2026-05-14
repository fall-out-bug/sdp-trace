package main

import (
	"io"
)

func writeOptionalPRReviewSummary(path, summary string, stderr io.Writer) (int, bool) {
	if path == "" {
		// Summary output is optional; no path means there is no write attempt and
		// no additional failure state.
		return 0, true
	}
	return writePRReviewSummaryFile(path, summary, stderr)
}
