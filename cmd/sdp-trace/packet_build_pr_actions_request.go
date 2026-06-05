package main

import (
	"net/http"
	"strings"
)

func githubActionsArtifactsRequest(ctx githubActionsArtifactContext) (*http.Request, error) {
	url := strings.TrimRight(ctx.apiURL, "/") + "/repos/" + ctx.repo + "/actions/runs/" + ctx.runID + "/artifacts"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Keep the GitHub API media type explicit so artifact evidence fetches are
	// stable across server defaults.
	req.Header.Set("Accept", "application/vnd.github+json")
	// Never attach a GitHub token to the loopback HTTP test API path.
	if auth := githubActionsArtifactsAuthorization(ctx.apiURL, ctx.token); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	return req, nil
}
