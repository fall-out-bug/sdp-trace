package main

import (
	"fmt"
	"os"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func completePRInputFromOptions(opts *flagSet, source string, input *packet.GitHubPREvidenceInput) error {
	if err := readOptionalPREvidence(opts, input); err != nil {
		return err
	}
	// Live GitHub hydration is skipped for fixture mode so local replay remains
	// hermetic.
	if err := hydrateGitHubActionsEvidence(source, opts.stringValue("github-api-url"), input, os.Getenv); err != nil {
		return err
	}
	route, err := readOptionalPRRoute(opts.stringValue("route-manifest"))
	if err != nil {
		return fmt.Errorf("read route manifest: %w", err)
	}
	// Route manifests are optional enrichment; an empty manifest leaves route
	// rows to validate as missing or cannot_verify.
	applyPRRoute(input, route)
	return nil
}
