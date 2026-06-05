package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func writePRReviewRunOutput(stdout io.Writer, runs prreview.RunSet, preview *prreview.RunPreview) {
	if preview != nil {
		// Preview mode reports planned runner invocations without implying that
		// any review evidence has been produced.
		writeIndentedPayload(stdout, preview)
		return
	}
	writeIndentedPayload(stdout, runs)
}
