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
