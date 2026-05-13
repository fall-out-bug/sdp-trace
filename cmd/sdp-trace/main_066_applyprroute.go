package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func applyPRRoute(input *packet.GitHubPREvidenceInput, route packet.GitHubPREvidenceInput) {
	// Route manifests overwrite only route/review fields, leaving PR identity and
	// CI evidence anchored to the selected event source.
	input.AgentRouteRefs = route.AgentRouteRefs
	input.AgentRouteComponents = route.AgentRouteComponents
	input.AgentRouteDigest = route.AgentRouteDigest
	input.AgentRouteEvidenceKind = route.AgentRouteEvidenceKind
	input.PromptBoundary = route.PromptBoundary
	input.IntegrationActions = route.IntegrationActions
	input.Reviews = route.Reviews
}
