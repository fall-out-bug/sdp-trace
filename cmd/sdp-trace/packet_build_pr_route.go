package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func readOptionalPRRoute(path string) (packet.GitHubPREvidenceInput, error) {
	var route packet.GitHubPREvidenceInput
	return route, readOptionalJSON(path, &route)
}

func applyPRRoute(input *packet.GitHubPREvidenceInput, route packet.GitHubPREvidenceInput) {
	// Route manifests overwrite only routing, prompt-boundary, integration-action,
	// and review fields, leaving PR identity and CI evidence anchored to the
	// selected event source.
	input.AgentRouteRefs = route.AgentRouteRefs
	input.AgentRouteComponents = route.AgentRouteComponents
	input.AgentRouteDigest = route.AgentRouteDigest
	input.AgentRouteEvidenceKind = route.AgentRouteEvidenceKind
	input.PromptBoundary = route.PromptBoundary
	input.IntegrationActions = route.IntegrationActions
	input.Reviews = route.Reviews
}
