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
	return completePRInputRoute(opts, input)
}

func completePRInputRoute(opts *flagSet, input *packet.GitHubPREvidenceInput) error {
	route, err := readOptionalPRRoute(opts.stringValue("route-manifest"))
	if err != nil {
		return fmt.Errorf("read route manifest: %w", err)
	}
	// Route manifests are optional enrichment; an empty manifest leaves route
	// rows to validate as missing or cannot_verify.
	applyPRRoute(input, route)
	if err := packet.ValidateGitHubPREvidenceInputResolvers(*input); err != nil {
		return fmt.Errorf("validate github evidence resolvers: %w", err)
	}
	return nil
}

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
