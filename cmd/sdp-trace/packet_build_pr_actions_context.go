package main

import (
	"errors"
	"strings"
)

func newGitHubActionsArtifactContext(apiURLFlag string, getenv func(string) string) (githubActionsArtifactContext, error) {
	apiURL, err := githubAPIURL(apiURLFlag, getenv)
	if err != nil {
		return githubActionsArtifactContext{}, err
	}
	// Context captures only the repository, run id, token, and validated API
	// endpoint needed for artifact listing.
	ctx := githubActionsArtifactContext{
		repo:   getenv("GITHUB_REPOSITORY"),
		runID:  getenv("GITHUB_RUN_ID"),
		token:  githubToken(getenv),
		apiURL: apiURL,
	}
	if err := validateGitHubActionsArtifactContext(ctx); err != nil {
		// Repository/run/token validation keeps failed discovery local and
		// deterministic before any artifact API request is built.
		return githubActionsArtifactContext{}, err
	}
	return ctx, nil
}

func validateGitHubActionsArtifactContext(ctx githubActionsArtifactContext) error {
	if missingGitHubArtifactIdentity(ctx) {
		// Repository and run id bind artifact discovery to the current PR run.
		return errors.New("missing GITHUB_REPOSITORY or GITHUB_RUN_ID for GitHub Actions artifact discovery")
	}
	if strings.TrimSpace(ctx.token) == "" {
		// Without a token, the command cannot prove artifact availability from the
		// configured run.
		return errors.New("missing GITHUB_TOKEN or GH_TOKEN for GitHub Actions artifact discovery")
	}
	return nil
}

func missingGitHubArtifactIdentity(ctx githubActionsArtifactContext) bool {
	return strings.TrimSpace(ctx.repo) == "" || strings.TrimSpace(ctx.runID) == ""
}
