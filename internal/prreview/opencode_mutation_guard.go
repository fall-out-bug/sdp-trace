package prreview

// applyOpenCodeMutationCheck invalidates OpenCode results when the runner
// mutates the working tree after a baseline was captured.
func applyOpenCodeMutationCheck(result *ReviewerResult, role ReviewRole, workDir string, baseline *workingTreeBaseline) {
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

func needsOpenCodeMutationCheck(role ReviewRole, baseline *workingTreeBaseline) bool {
	return role.Runner == RunnerOpenCode && baseline != nil
}

func markBaselineCannotVerify(result *ReviewerResult) {
	result.Status = StatusCannotVerify
	result.Reason = "working_tree_baseline_cannot_verify"
}

func baselineChanged(after, before *workingTreeBaseline) bool {
	return after.Digest != before.Digest || after.Count != before.Count
}
