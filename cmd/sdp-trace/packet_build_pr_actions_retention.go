package main

import (
	"fmt"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func retainedGitHubArtifacts(payload githubActionsArtifactPayload, ctx githubActionsArtifactContext) []packet.GitHubArtifact {
	artifacts := []packet.GitHubArtifact{}
	for _, artifact := range payload.Artifacts {
		if artifact.Expired {
			// Expired artifacts are not retained evidence.
			continue
		}
		artifacts = append(artifacts, packet.GitHubArtifact{
			Name:         artifact.Name,
			Resolver:     githubArtifactResolver(artifact, ctx),
			RetainedForm: "external_ref",
			ExpiresAt:    artifact.ExpiresAt,
		})
	}
	// Retained artifacts are external references only; packet validation decides
	// whether they satisfy the required evidence rows.
	return artifacts
}

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
