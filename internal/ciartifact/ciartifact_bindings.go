package ciartifact

func bindingSummary(families []FamilyObservation) BindingSummary {
	// Binding summaries keep producer scope and access state separate because each
	// can fail independently.

	sourceRun := BindingNotAssessed
	producer := BindingNotAssessed
	for _, family := range families {
		if !family.Required {

			continue
		}
		sourceRun = worseBinding(sourceRun, family.BindingState)
		producer = worseBinding(producer, producerBindingState(family))
	}
	return BindingSummary{SourceBindingState: sourceRun, RunBindingState: sourceRun, ProducerBindingState: producer}
}

func producerBindingState(family FamilyObservation) string {
	// producerBindingState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if family.RequiredProducer == ProducerCIUploaded && family.ProducerScope != ProducerCIUploaded {

		return BindingMismatch
	}
	return BindingMatched
}

func worseBinding(current, candidate string) string {
	// worseBinding keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if bindingRank(candidate) > bindingRank(current) {

		return candidate
	}
	return current
}

func bindingRank(state string) int {
	// bindingRank keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if rank, ok := bindingRanks[state]; ok {
		return rank
	}

	return 1
}

var bindingRanks = map[string]int{
	BindingMismatch:     5,
	BindingAbsent:       4,
	BindingUnverifiable: 3,
	BindingMatched:      2,
}
