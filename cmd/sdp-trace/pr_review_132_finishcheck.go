package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func finishPRReviewCheck(outDir string, packet prreview.Packet, profile prreview.ReviewProfile, runs prreview.RunSet, preview *prreview.RunPreview, stdout, stderr io.Writer) int {
	if writePRReviewCheckPreview(stdout, preview) {
		return 0
	}
	// Non-preview mode persists artifacts before printing the human summary, so
	// the summary never outruns the machine-readable evidence.
	ledger, validation, code, ok := writePRReviewCheckArtifacts(outDir, packet, profile, runs, stderr)
	if !ok {
		return code
	}
	fmt.Fprint(stdout, prreview.Summarize(validation, ledger))
	return reviewValidationExit(validation)
}
