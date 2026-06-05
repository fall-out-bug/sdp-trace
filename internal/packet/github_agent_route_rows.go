package packet

import "strings"

// GitHub route rows distinguish missing route evidence from prompt-boundary
// failures. A failed prompt boundary blocks route proof before route refs can
// be treated as partial evidence.
func githubAgentRouteRow(input GitHubPREvidenceInput) Row {
	classification := ClassifyPromptBoundary(input.PromptBoundary)
	if row, ok := promptBoundaryRouteFailureRow(input.RequirePromptBoundary, classification); ok {
		return row
	}

	if len(input.AgentRouteRefs) > 0 {
		return githubAgentRouteRefsRow(input.RequirePromptBoundary, classification)
	}
	return githubRow("PC-AGENT-ROUTE", StateNotAssessed, "Agent route evidence was not provided.", nil, "missing OpenCode/GSD observation ref")
}

func githubAgentRouteRefsRow(requirePromptBoundary bool, classification PromptBoundaryClassification) Row {
	if requirePromptBoundary && classification.RouteProofEffect == StatePartial {
		return githubRow("PC-AGENT-ROUTE", StatePartial, "Agent route refs and digest-only prompt boundary are retained.", []string{"agent:route", "prompt:boundary"}, "prompt text is unavailable; digest-only boundary supports partial route proof")
	}
	return githubRow("PC-AGENT-ROUTE", StatePartial, "Agent route refs are retained.", []string{"agent:route"}, "route refs are input refs, not a complete observed delegation chain")
}

func githubPromptBoundaryRouteFailRow(classification PromptBoundaryClassification) Row {
	return githubRow("PC-AGENT-ROUTE", StateFail, "Developer prompt contains recorder duties.", []string{"prompt:boundary"}, strings.Join(classification.Reasons, "; "))
}

func githubPromptBoundaryRouteCannotVerifyRow(classification PromptBoundaryClassification) Row {
	return githubRow("PC-AGENT-ROUTE", StateCannotVerify, "Prompt boundary evidence cannot verify developer-route independence.", []string{"prompt:boundary"}, strings.Join(classification.Reasons, "; "))
}
