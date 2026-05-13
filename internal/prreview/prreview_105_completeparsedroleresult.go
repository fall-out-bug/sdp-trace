package prreview

func completeParsedRoleResult(parsed ReviewerResult, role ReviewRole, workDir string, baseline *workingTreeBaseline) ReviewerResult {
	// completeParsedRoleResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	applyOpenCodeMutationCheck(&parsed, role, workDir, baseline)
	return parsed
}
