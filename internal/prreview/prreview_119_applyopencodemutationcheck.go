package prreview

func applyOpenCodeMutationCheck(result *ReviewerResult, role ReviewRole, workDir string, baseline *workingTreeBaseline) {
	// applyOpenCodeMutationCheck keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if !needsOpenCodeMutationCheck(role, baseline) {
		return
	}
	after, err := captureWorkingTreeBaseline(workDir)
	if err != nil {
		markBaselineCannotVerify(result)
		return
	}
	if baselineChanged(after, baseline) {
		result.Status = StatusCannotVerify
		result.Reason = "mutation_detected"
	}
}
