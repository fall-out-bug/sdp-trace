package prreview

func completeUnparsedRoleErrorOrEmpty(result ReviewerResult, output []byte, runErr error) (ReviewerResult, bool) {
	// completeUnparsedRoleErrorOrEmpty keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	switch {
	case applyRunnerError(&result, runErr) != nil:
		return result, true
	case emptyReviewerOutput(output):
		result.Status = StatusEmptyOutput
		result.Reason = "runner_empty_output"
		return result, true
	default:
		return result, false
	}
}
