package prreview

// prepareOpenCodeBaseline captures a clean working-tree baseline before an
// OpenCode runner can execute.
func prepareOpenCodeBaseline(result *ReviewerResult, role ReviewRole, workDir string) (*workingTreeBaseline, bool, error) {
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

// openCodeBaselineReady enforces clean-required mode before command execution.
func openCodeBaselineReady(result *ReviewerResult, role ReviewRole, baseline *workingTreeBaseline) (*workingTreeBaseline, bool, error) {
	if !openCodeBaselineClean(result, role, baseline) {
		return nil, false, nil
	}

	return baseline, commandConfigured(result, role), nil
}

// openCodeReadOnlyReady requires explicit read-only enforcement before running
// OpenCode.
func openCodeReadOnlyReady(result *ReviewerResult, role ReviewRole) bool {
	if !role.ReadOnlyEnforced {
		markOpenCodeReadOnlyMissing(result)
		return false
	}
	if err := attachPromptRef(result, role); err != nil {
		return false
	}
	return true
}

// openCodeBaselineClean rejects dirty clean-required worktrees as not_assessed
// rather than treating them as a reviewer result.
func openCodeBaselineClean(result *ReviewerResult, role ReviewRole, baseline *workingTreeBaseline) bool {
	mode := defaultString(role.WorkingTreeMode, "clean_required")
	if mode != "clean_required" || baseline.Count == 0 {
		return true
	}
	result.Status = StatusNotAssessed
	result.Reason = "working_tree_dirty"
	return false
}

// markOpenCodeReadOnlyMissing records the missing read-only precondition.
func markOpenCodeReadOnlyMissing(result *ReviewerResult) {
	result.Reason = "opencode_read_only_not_enforced"
}
