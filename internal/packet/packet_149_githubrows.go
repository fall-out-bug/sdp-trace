package packet

func githubRows(input GitHubPREvidenceInput) []Row {
	// githubRows keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	change := githubChangeRow(input)
	mutation := githubMutationRow(input)

	return []Row{
		change,
		githubInitiatorRow(input),
		githubAgentRouteRow(input),
		mutation,
		githubVerificationRow(input),
		githubReviewRow(input),
		githubRow("PC-AUTHORITY", StateNotAssessed, "Authority was not assessed for this generated GitHub input.", nil, "authority profile was not provided"),
		githubRow("PC-THEATER", StatePass, "No P0 theater finding triggered by the minimal GitHub input builder.", []string{"theater:builder"}, ""),
		githubRow("PC-ATTESTATION", StateNotAssessed, "Signed or external attestation was not assessed.", nil, "signed trust inputs were not provided"),
		githubRow("PC-DECISION", StateNotAssessed, "Default decision owner placeholders are recorded.", []string{"decision:owners"}, "decision owners are placeholders, not bound approval or ownership evidence"),
		githubRow("PC-RESIDUAL-GAPS", StatePartial, "Non-pass rows remain explicit in residual gaps.", []string{"gap:generated"}, "generated packet contains explicit non-pass rows"),
	}
}
