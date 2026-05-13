package prreview

func completeRoleResult(result ReviewerResult, role ReviewRole, packet Packet, workDir string, baseline *workingTreeBaseline, output []byte, timedOut bool, runErr error) ReviewerResult {
	// completeRoleResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if completed, ok := completeUnparsedRoleResult(result, output, timedOut, runErr); ok {
		return completed
	}
	parsed, err := parseReviewerOutput(result, role, packet, output)
	return completeParsedRoleResult(markRoleParseFailure(parsed, err), role, workDir, baseline)
}
