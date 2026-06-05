package prreview

// completeRoleResult maps runner output into a reviewer result while keeping
// unparsed runner failures distinct from parsed reviewer findings.
func completeRoleResult(result ReviewerResult, role ReviewRole, packet Packet, workDir string, baseline *workingTreeBaseline, output []byte, timedOut bool, runErr error) ReviewerResult {
	if completed, ok := completeUnparsedRoleResult(result, output, timedOut, runErr); ok {
		return completed
	}
	parsed, err := parseReviewerOutput(result, role, packet, output)
	return completeParsedRoleResult(markRoleParseFailure(parsed, err), role, workDir, baseline)
}

// completeParsedRoleResult applies post-parse trust checks that need parsed
// reviewer status and role runner context.
func completeParsedRoleResult(parsed ReviewerResult, role ReviewRole, workDir string, baseline *workingTreeBaseline) ReviewerResult {
	applyOpenCodeMutationCheck(&parsed, role, workDir, baseline)
	return parsed
}

// markRoleParseFailure preserves parsed defaults while making malformed runner
// output non-authoritative.
func markRoleParseFailure(result ReviewerResult, err error) ReviewerResult {
	if err == nil {
		return result
	}
	result.Status = StatusParseFailed
	result.Reason = "runner_output_parse_failed"
	return result
}
