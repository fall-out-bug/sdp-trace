package prreview

func openCodeBaselineReady(result *ReviewerResult, role ReviewRole, baseline *workingTreeBaseline) (*workingTreeBaseline, bool, error) {
	// openCodeBaselineReady keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if !openCodeBaselineClean(result, role, baseline) {
		return nil, false, nil
	}

	return baseline, commandConfigured(result, role), nil
}
