package prreview

func prepareOpenCodeBaseline(result *ReviewerResult, role ReviewRole, workDir string) (*workingTreeBaseline, bool, error) {
	// prepareOpenCodeBaseline keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if !openCodeReadOnlyReady(result, role) {
		return nil, false, nil
	}
	baseline, err := captureWorkingTreeBaseline(workDir)
	if err != nil {
		result.Status = StatusCannotVerify
		result.Reason = "working_tree_baseline_cannot_verify"
		return nil, false, nil
	}
	return openCodeBaselineReady(result, role, baseline)
}
