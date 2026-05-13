package prreview

func completeUnparsedRoleResult(result ReviewerResult, output []byte, timedOut bool, runErr error) (ReviewerResult, bool) {
	// completeUnparsedRoleResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if timedOut {
		result.Status = StatusTimedOut
		result.Reason = "runner_timed_out"
		return result, true
	}
	return completeUnparsedRoleErrorOrEmpty(result, output, runErr)
}
