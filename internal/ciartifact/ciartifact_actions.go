package ciartifact

func nextActions(families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult, identityCannotVerify bool) []string {
	// Next actions name the smallest missing or unsafe evidence boundary needed to
	// move the CI-artifact observation forward.
	set := map[string]bool{}

	addFamilyActions(set, families)
	addConditionalAction(set, resultNeedsAction(index.Result), "Regenerate or supply a verifier-readable artifact index.")
	addConditionalAction(set, resultNeedsAction(safety.State), "Use the recorded safety ruleset id to remove unsafe artifact output before rerun.")
	addConditionalAction(set, identityCannotVerify, "Provide safe source and run identity metadata before using this observation as CI-backed proof.")
	return sortedKeys(set)
}

func addFamilyActions(set map[string]bool, families []FamilyObservation) {
	// addFamilyActions keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, family := range families {

		addConditionalAction(set, family.NextAction != "", family.NextAction)
	}
}

func addConditionalAction(set map[string]bool, include bool, action string) {
	if include {
		// The map deduplicates identical operator actions from multiple evidence rows.
		set[action] = true
	}
}

func resultNeedsAction(state string) bool {
	return state == StateFail || state == StateCannotVerify
}
