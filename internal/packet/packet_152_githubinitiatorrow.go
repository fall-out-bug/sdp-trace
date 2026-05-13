package packet

func githubInitiatorRow(input GitHubPREvidenceInput) Row {
	// githubInitiatorRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if input.PR.BodyRef != "" {

		return githubRow("PC-INITIATOR", StatePartial, "PR body task source is retained.", []string{"github:pr-body"}, "PR body is weaker than a dedicated issue binding")
	}
	return githubRow("PC-INITIATOR", StateNotAssessed, "No task or initiator evidence was provided.", nil, "missing PR body, issue, or retained task artifact")
}
