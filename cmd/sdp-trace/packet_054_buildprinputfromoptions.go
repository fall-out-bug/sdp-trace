package main

import (
	"os"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func buildPRInputFromOptions(opts *flagSet) (packet.GitHubPREvidenceInput, error) {
	source, event, err := loadPRInputSourceEvent(opts)
	if err != nil {
		return packet.GitHubPREvidenceInput{}, err
	}
	// Event metadata seeds the input before optional local and live evidence is
	// layered in.
	input := githubPRInputFromEvent(event, source, os.Getenv)
	if err := completePRInputFromOptions(opts, source, &input); err != nil {
		// Optional evidence failures still invalidate the whole packet input,
		// because partial PR packets can overstate route or CI readiness.
		return packet.GitHubPREvidenceInput{}, err
	}
	return input, nil
}
