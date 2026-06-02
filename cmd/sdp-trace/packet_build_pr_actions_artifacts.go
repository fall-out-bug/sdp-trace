package main

import (
	"errors"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

type githubActionsArtifactContext struct {
	repo   string
	runID  string
	token  string
	apiURL string
}

type githubActionsArtifactPayload struct {
	Artifacts []githubActionsArtifact `json:"artifacts"`
}

type githubActionsArtifact struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Expired   bool   `json:"expired"`
	ExpiresAt string `json:"expires_at"`
	URL       string `json:"archive_download_url"`
}

func githubActionsArtifacts(apiURLFlag string, getenv func(string) string) ([]packet.GitHubArtifact, error) {
	ctx, err := newGitHubActionsArtifactContext(apiURLFlag, getenv)
	if err != nil {
		return nil, err
	}
	// Artifact discovery is live GitHub evidence and therefore requires a fully
	// validated request context before the network call.
	payload, err := fetchGitHubActionsArtifacts(ctx)
	if err != nil {
		return nil, err
	}
	artifacts := retainedGitHubArtifacts(payload, ctx)
	if len(artifacts) == 0 {
		// An empty retained set means there is no durable artifact evidence for
		// the packet to cite.
		return nil, errors.New("GitHub Actions artifact discovery returned no retained artifacts")
	}
	return artifacts, nil
}
