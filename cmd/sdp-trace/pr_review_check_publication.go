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

func writePRReviewCheckPreview(stdout io.Writer, preview *prreview.RunPreview) bool {
	if preview == nil {
		return false
	}
	// Preview output is intentionally terminal-only planning data, not persisted
	// review evidence.
	writeIndentedPayload(stdout, preview)
	return true
}

func reviewValidationExit(validation prreview.Validation) int {
	if reviewValidationExitCode(validation) != 0 {
		// Synthesis validation failures are trust gaps, not usage errors.
		return exitCannotVerify
	}
	return 0
}
