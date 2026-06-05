package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func successfulHTTPStatus(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

func decodeGitHubActionsArtifacts(reader io.Reader) (githubActionsArtifactPayload, error) {
	var payload githubActionsArtifactPayload
	if err := json.NewDecoder(reader).Decode(&payload); err != nil {
		return githubActionsArtifactPayload{}, fmt.Errorf("decode GitHub Actions artifacts: %w", err)
	}
	// Payload schema validation is minimal here; packet row validation decides
	// whether retained artifact refs are sufficient evidence.
	return payload, nil
}
