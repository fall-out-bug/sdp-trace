package packet

func promptBoundaryRouteCannotVerifyRow(classification PromptBoundaryClassification) (Row, bool) {
	// promptBoundaryRouteCannotVerifyRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if classification.RouteProofEffect != StateCannotVerify {
		return Row{}, false
	}

	return githubPromptBoundaryRouteCannotVerifyRow(classification), true
}
