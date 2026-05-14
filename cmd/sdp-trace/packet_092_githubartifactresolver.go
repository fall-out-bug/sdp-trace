package main

import (
	"fmt"
	"strings"
)

func githubArtifactResolver(artifact githubActionsArtifact, ctx githubActionsArtifactContext) string {
	if artifact.URL != "" {
		return artifact.URL
	}
	if artifact.ID == 0 {
		// Without a resolver URL or id, downstream packet rows cannot cite the
		// artifact.
		return ""
	}
	return strings.TrimRight(ctx.apiURL, "/") + "/repos/" + ctx.repo + "/actions/artifacts/" + fmt.Sprint(artifact.ID) + "/zip"
}
