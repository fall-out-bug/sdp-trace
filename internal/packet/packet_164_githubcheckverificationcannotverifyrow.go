package packet

func githubCheckVerificationCannotVerifyRow(input GitHubPREvidenceInput) (Row, bool) {
	// githubCheckVerificationCannotVerifyRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if len(input.Checks) == 0 {

		return githubRow("PC-VERIFICATION", StateCannotVerify, "No GitHub check evidence was provided.", nil, "missing GitHub check or workflow run evidence"), true
	}
	if missingRequiredWorkflowRunID(input) {

		return githubRow("PC-VERIFICATION", StateCannotVerify, "No current workflow run id was provided.", []string{"github:check"}, "missing workflow run id for CI-owned packet generation"), true
	}
	return Row{}, false
}
