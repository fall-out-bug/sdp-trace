package main

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func readOptionalPREvidence(opts *flagSet, input *packet.GitHubPREvidenceInput) error {
	if err := readOptionalJSON(opts.stringValue("checks-json"), &input.Checks); err != nil {
		return fmt.Errorf("read checks json: %w", err)
	}
	// Artifacts can be provided by file or discovered from GitHub, but malformed
	// local artifact JSON is never ignored.
	if err := readOptionalJSON(opts.stringValue("artifacts-json"), &input.Artifacts); err != nil {
		return fmt.Errorf("read artifacts json: %w", err)
	}
	return nil
}
