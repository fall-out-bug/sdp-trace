package packet

func githubVerificationCannotVerifyRow(input GitHubPREvidenceInput, classification PromptBoundaryClassification) (Row, bool) {
	// githubVerificationCannotVerifyRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	if row, ok := githubPromptBoundaryVerificationCannotVerifyRow(input.RequirePromptBoundary, classification); ok {
		return row, true
	}
	return githubCheckVerificationCannotVerifyRow(input)
}
