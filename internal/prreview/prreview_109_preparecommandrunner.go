package prreview

func prepareCommandRunner(result *ReviewerResult, role ReviewRole) (*workingTreeBaseline, bool, error) {
	// prepareCommandRunner keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if err := attachPromptRef(result, role); err != nil {
		return nil, false, nil
	}

	return nil, commandConfigured(result, role), nil
}
