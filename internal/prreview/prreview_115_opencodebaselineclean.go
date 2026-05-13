package prreview

func openCodeBaselineClean(result *ReviewerResult, role ReviewRole, baseline *workingTreeBaseline) bool {
	// openCodeBaselineClean keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	mode := defaultString(role.WorkingTreeMode, "clean_required")
	if mode != "clean_required" || baseline.Count == 0 {
		return true
	}
	result.Status = StatusNotAssessed
	result.Reason = "working_tree_dirty"
	return false
}
