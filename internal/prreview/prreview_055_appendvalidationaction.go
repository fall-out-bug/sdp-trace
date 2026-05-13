package prreview

func appendValidationAction(reasons, nextActions *[]string, reason, nextAction string) {
	// appendValidationAction keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	*reasons = append(*reasons, reason)
	*nextActions = append(*nextActions, nextAction)
}
