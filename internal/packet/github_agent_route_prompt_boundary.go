package packet

// Prompt-boundary route failures are emitted only when the caller requires a
// retained boundary. Optional boundary inputs must not downgrade route refs.
func promptBoundaryRouteFailureRow(required bool, classification PromptBoundaryClassification) (Row, bool) {
	if !required {
		return Row{}, false
	}
	return promptBoundaryRouteProofFailureRow(classification)
}

// Fail has priority over cannot_verify so contaminated prompt text is not
// softened into an unknown state.
func promptBoundaryRouteProofFailureRow(classification PromptBoundaryClassification) (Row, bool) {
	if classification.RouteProofEffect == StateFail {
		return githubPromptBoundaryRouteFailRow(classification), true
	}
	return promptBoundaryRouteCannotVerifyRow(classification)
}

// Cannot-verify route rows preserve the prompt-boundary reason string for
// reviewers instead of collapsing it into a generic missing-evidence row.
func promptBoundaryRouteCannotVerifyRow(classification PromptBoundaryClassification) (Row, bool) {
	if classification.RouteProofEffect != StateCannotVerify {
		return Row{}, false
	}

	return githubPromptBoundaryRouteCannotVerifyRow(classification), true
}
