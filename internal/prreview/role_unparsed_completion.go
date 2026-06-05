package prreview

import (
	"strings"
)

// completeUnparsedRoleResult handles timeout and runner-output absence before
// attempting to parse structured reviewer JSON.
func completeUnparsedRoleResult(result ReviewerResult, output []byte, timedOut bool, runErr error) (ReviewerResult, bool) {
	if timedOut {
		result.Status = StatusTimedOut
		result.Reason = "runner_timed_out"
		return result, true
	}
	return completeUnparsedRoleErrorOrEmpty(result, output, runErr)
}

// completeUnparsedRoleErrorOrEmpty classifies process errors and empty output
// as runner states rather than reviewer verdicts.
func completeUnparsedRoleErrorOrEmpty(result ReviewerResult, output []byte, runErr error) (ReviewerResult, bool) {
	if applyRunnerError(&result, runErr) != nil {
		return result, true
	}
	if emptyReviewerOutput(output) {
		result.Status = StatusEmptyOutput
		result.Reason = "runner_empty_output"
		return result, true
	}
	return result, false
}

// emptyReviewerOutput treats whitespace-only output as no reviewer evidence.
func emptyReviewerOutput(output []byte) bool {
	return len(strings.TrimSpace(string(output))) == 0
}
