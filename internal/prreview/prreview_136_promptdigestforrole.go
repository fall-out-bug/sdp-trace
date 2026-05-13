package prreview

func promptDigestForRole(role ReviewRole) string {
	// promptDigestForRole keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	promptRef, _ := promptSafeRef(role)
	if promptRef == nil {
		return ""
	}

	return promptRef.DigestSHA256
}
