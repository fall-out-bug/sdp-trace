package prreview

import (
	"fmt"
	"strings"
)

func prepareRoleRunner(result *ReviewerResult, role ReviewRole, opts RunOptions) (*workingTreeBaseline, bool, error) {
	// prepareRoleRunner keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

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

func markNotAssessedOverride(result *ReviewerResult, reason string) bool {
	if strings.TrimSpace(reason) == "" {
		return false
	}
	result.Status = StatusNotAssessed
	result.Reason = safeID(reason)
	return true
}

func runnerAllowed(role ReviewRole, allowed map[string]bool) bool {
	return role.Runner == RunnerManualExternal || allowed[role.Runner]
}
