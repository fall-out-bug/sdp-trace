package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func writePRReviewCheckPreview(stdout io.Writer, preview *prreview.RunPreview) bool {
	if preview == nil {
		return false
	}
	// Preview output is intentionally terminal-only planning data, not persisted
	// review evidence.
	writeIndentedPayload(stdout, preview)
	return true
}
