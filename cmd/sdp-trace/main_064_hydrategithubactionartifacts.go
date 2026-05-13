package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func hydrateGitHubActionArtifacts(apiURL string, input *packet.GitHubPREvidenceInput, getenv func(string) string) error {
	if len(input.Artifacts) != 0 {
		// Explicit artifact JSON wins over live discovery for replayability.
		return nil
	}
	artifacts, err := githubActionsArtifacts(apiURL, getenv)
	if err != nil {
		return err
	}
	input.Artifacts = artifacts
	return nil
}
