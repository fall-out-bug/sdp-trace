package prreview

func markRoleParseFailure(result ReviewerResult, err error) ReviewerResult {
	// markRoleParseFailure keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if err == nil {
		return result
	}
	result.Status = StatusParseFailed
	result.Reason = "runner_output_parse_failed"
	return result
}
