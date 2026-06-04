package prreview

import (
	"fmt"
	"strings"
)

// prepareRoleRunner applies caller policy before any external runner is
// invoked.
func prepareRoleRunner(result *ReviewerResult, role ReviewRole, opts RunOptions) (*workingTreeBaseline, bool, error) {
	if markNotAssessedOverride(result, opts.NotAssessedReason) {
		return nil, false, nil
	}
	if !runnerAllowed(role, opts.AllowedRunners) {
		return nil, false, fmt.Errorf("runner_not_allowed: %s", role.Runner)
	}
	if role.Runner == RunnerOpenCode {
		return prepareOpenCodeBaseline(result, role, opts.WorkDir)
	}
	return prepareCommandRunner(result, role)
}

// markNotAssessedOverride records an explicit caller-provided reason without
// invoking a model or runner.
func markNotAssessedOverride(result *ReviewerResult, reason string) bool {
	if strings.TrimSpace(reason) == "" {
		return false
	}
	result.Status = StatusNotAssessed
	result.Reason = safeID(reason)
	return true
}

// runnerAllowed keeps executable runners opt-in while manual imports remain
// local evidence.
func runnerAllowed(role ReviewRole, allowed map[string]bool) bool {
	return role.Runner == RunnerManualExternal || allowed[role.Runner]
}
