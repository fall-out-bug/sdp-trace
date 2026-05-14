package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func hydrateGitHubActionsEvidence(source string, apiURL string, input *packet.GitHubPREvidenceInput, getenv func(string) string) error {
	if source != "github-actions" {
		// Fixture mode must not make network calls.
		return nil
	}
	if err := hydrateGitHubActionArtifacts(apiURL, input, getenv); err != nil {
		return err
	}
	return nil
}
