package main

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func writePRReviewCheckArtifacts(outDir string, packet prreview.Packet, profile prreview.ReviewProfile, runs prreview.RunSet, stderr io.Writer) (prreview.Ledger, prreview.Validation, int, bool) {
	if !writePRReviewJSON(filepath.Join(outDir, "runs", "results.json"), runs, stderr) {
		// Run-set persistence is the first durable artifact in non-preview mode.
		return prreview.Ledger{}, prreview.Validation{}, 1, false
	}
	// Ledger and validation are derived after run persistence so artifact paths
	// and in-memory validation stay in the same review cycle.
	ledger := prreview.SynthesizeLedger(packet, runs, nil)
	validation := prreview.Validate(packet, profile, runs, ledger)
	if !writePRReviewJSON(filepath.Join(outDir, "ledger.json"), ledger, stderr) {
		return prreview.Ledger{}, prreview.Validation{}, 1, false
	}
	// Validation is written after ledger so readers never see a validation file
	// whose cited ledger is missing.
	if !writePRReviewJSON(filepath.Join(outDir, "validation.json"), validation, stderr) {
		return prreview.Ledger{}, prreview.Validation{}, 1, false
	}
	return ledger, validation, 0, true
}

func writePRReviewJSON(path string, value any, stderr io.Writer) bool {
	if err := prreview.WriteJSON(path, value); err != nil {
		// Artifact write failure means the review evidence cannot be cited later.
		fmt.Fprintln(stderr, err)
		return false
	}
	return true
}
