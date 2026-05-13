package ciartifact

func orderedRequirements(reqs map[string]FamilyRequirement) []FamilyRequirement {
	// orderedRequirements keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	out := make([]FamilyRequirement, 0, len(reqs))
	for _, family := range familyOrder {
		if req, ok := reqs[family]; ok {
			out = append(out, req)
		}
	}
	return out
}
