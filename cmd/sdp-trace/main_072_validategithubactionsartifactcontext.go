package main

import (
	"errors"
	"strings"
)

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
