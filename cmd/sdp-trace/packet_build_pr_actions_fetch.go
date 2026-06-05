package main

import (
	"fmt"
	"net/http"
)

func fetchGitHubActionsArtifacts(ctx githubActionsArtifactContext) (githubActionsArtifactPayload, error) {
	req, err := githubActionsArtifactsRequest(ctx)
	if err != nil {
		// Request construction failures happen before credentials leave the
		// process.
		return githubActionsArtifactPayload{}, err
	}
	// This is the only live network fetch in packet build-pr artifact hydration.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return githubActionsArtifactPayload{}, fmt.Errorf("list GitHub Actions artifacts: %w", err)
	}
	defer resp.Body.Close()
	if !successfulHTTPStatus(resp.StatusCode) {
		// Non-2xx responses mean the retained artifact set cannot be verified.
		return githubActionsArtifactPayload{}, fmt.Errorf("list GitHub Actions artifacts: HTTP %d", resp.StatusCode)
	}
	// Decoding is the handoff from live GitHub response bytes to packet evidence
	// candidates; retained-artifact filtering happens after this step.
	return decodeGitHubActionsArtifacts(resp.Body)
}
