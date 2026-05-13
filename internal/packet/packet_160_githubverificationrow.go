package packet

func githubVerificationRow(input GitHubPREvidenceInput) Row {
	// githubVerificationRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	classification := ClassifyPromptBoundary(input.PromptBoundary)
	if row, ok := githubVerificationCannotVerifyRow(input, classification); ok {
		return row
	}

	if !checksHaveRetainedArtifactRefs(input) {
		return githubRow("PC-VERIFICATION", StatePartial, "GitHub check evidence is retained without retained artifact binding.", []string{"github:check"}, "GitHub CI green is not verification pass without retained artifact evidence")
	}
	if !checksSucceeded(input.Checks) {
		return githubRow("PC-VERIFICATION", StatePartial, "GitHub checks include non-success conclusions.", []string{"github:check"}, "not all retained checks concluded success")
	}
	return githubVerificationPassRow(input)
}
