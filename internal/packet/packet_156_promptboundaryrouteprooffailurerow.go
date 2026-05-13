package packet

func promptBoundaryRouteProofFailureRow(classification PromptBoundaryClassification) (Row, bool) {
	// promptBoundaryRouteProofFailureRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if classification.RouteProofEffect == StateFail {

		return githubPromptBoundaryRouteFailRow(classification), true
	}
	return promptBoundaryRouteCannotVerifyRow(classification)
}
