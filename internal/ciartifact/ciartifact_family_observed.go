package ciartifact

func observedFamilies(inputs []FamilyInput) map[string]FamilyInput {
	// Observed families are collected from sanitized manifest data so unsafe labels
	// cannot become source references in the final assessment.
	observed := map[string]FamilyInput{}
	for _, input := range inputs {

		input.Family = canonicalFamily(input.Family)
		if !validFamily(input.Family) {
			continue
		}
		input.ProducerScope = safeProducerScope(input.ProducerScope)
		input.ArtifactAccessState = safeAccessState(input.ArtifactAccessState)
		input.BindingState = safeBindingState(input.BindingState)

		observed[input.Family] = input
	}
	return observed
}
