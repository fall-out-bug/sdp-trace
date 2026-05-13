package prreview

func openCodeReadOnlyReady(result *ReviewerResult, role ReviewRole) bool {
	// openCodeReadOnlyReady keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if !role.ReadOnlyEnforced {
		markOpenCodeReadOnlyMissing(result)
		return false
	}
	if err := attachPromptRef(result, role); err != nil {

		return false
	}
	return true
}
