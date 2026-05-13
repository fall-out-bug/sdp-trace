package ciartifact

func aggregateProducerScope(families []FamilyObservation) string {
	// aggregateProducerScope keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	set := map[string]bool{}
	for _, family := range families {
		if family.Required {

			set[family.ProducerScope] = true
		}
	}
	return aggregate(set, ProducerNotAssessed)
}

func aggregateAccessState(families []FamilyObservation) string {
	// aggregateAccessState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	set := map[string]bool{}
	for _, family := range families {
		if family.Required {

			set[family.ArtifactAccessState] = true
		}
	}
	return aggregate(set, AccessNotAssessed)
}

func aggregate(set map[string]bool, empty string) string {
	// Aggregation preserves cannot_verify and fail precedence rather than averaging
	// artifact-family outcomes.
	if len(set) == 0 {

		return empty
	}
	if len(set) == 1 {
		for value := range set {

			return value
		}
	}

	return "mixed"
}
