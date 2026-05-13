package packet

func githubAgentRouteRow(input GitHubPREvidenceInput) Row {
	// githubAgentRouteRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	classification := ClassifyPromptBoundary(input.PromptBoundary)
	if row, ok := promptBoundaryRouteFailureRow(input.RequirePromptBoundary, classification); ok {
		return row
	}

	if len(input.AgentRouteRefs) > 0 {
		return githubAgentRouteRefsRow(input.RequirePromptBoundary, classification)
	}
	return githubRow("PC-AGENT-ROUTE", StateNotAssessed, "Agent route evidence was not provided.", nil, "missing OpenCode/GSD observation ref")
}
