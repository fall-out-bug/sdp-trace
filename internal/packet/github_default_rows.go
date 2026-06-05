package packet

// Generated GitHub packets keep non-derived trust areas explicit. These rows
// prevent absent authority, attestation, decision, or residual-gap evidence
// from being mistaken for an implicit pass.
func githubAuthorityRow() Row {
	return githubRow("PC-AUTHORITY", StateNotAssessed, "Authority was not assessed for this generated GitHub input.", nil, "authority profile was not provided")
}

func githubTheaterRow() Row {
	return githubRow("PC-THEATER", StatePass, "No P0 theater finding triggered by the minimal GitHub input builder.", []string{"theater:builder"}, "")
}

func githubAttestationRow() Row {
	return githubRow("PC-ATTESTATION", StateNotAssessed, "Signed or external attestation was not assessed.", nil, "signed trust inputs were not provided")
}

func githubDecisionRow() Row {
	return githubRow("PC-DECISION", StateNotAssessed, "Default decision owner placeholders are recorded.", []string{"decision:owners"}, "decision owners are placeholders, not bound approval or ownership evidence")
}

func githubResidualGapsRow() Row {
	return githubRow("PC-RESIDUAL-GAPS", StatePartial, "Non-pass rows remain explicit in residual gaps.", []string{"gap:generated"}, "generated packet contains explicit non-pass rows")
}
