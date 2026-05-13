package main

import (
	"strings"
)

func missingGitHubArtifactIdentity(ctx githubActionsArtifactContext) bool {
	return strings.TrimSpace(ctx.repo) == "" || strings.TrimSpace(ctx.runID) == ""
}
