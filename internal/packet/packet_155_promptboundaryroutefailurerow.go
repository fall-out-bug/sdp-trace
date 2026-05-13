package packet

func promptBoundaryRouteFailureRow(required bool, classification PromptBoundaryClassification) (Row, bool) {
	// promptBoundaryRouteFailureRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if !required {

		return Row{}, false
	}
	return promptBoundaryRouteProofFailureRow(classification)
}
