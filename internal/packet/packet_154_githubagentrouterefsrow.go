package packet

func githubAgentRouteRefsRow(requirePromptBoundary bool, classification PromptBoundaryClassification) Row {
	// githubAgentRouteRefsRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	if requirePromptBoundary && classification.RouteProofEffect == StatePartial {
		return githubRow("PC-AGENT-ROUTE", StatePartial, "Agent route refs and digest-only prompt boundary are retained.", []string{"agent:route", "prompt:boundary"}, "prompt text is unavailable; digest-only boundary supports partial route proof")
	}
	return githubRow("PC-AGENT-ROUTE", StatePartial, "Agent route refs are retained.", []string{"agent:route"}, "route refs are input refs, not a complete observed delegation chain")
}
