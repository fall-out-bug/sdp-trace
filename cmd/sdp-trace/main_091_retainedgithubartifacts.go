package main

import (
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
