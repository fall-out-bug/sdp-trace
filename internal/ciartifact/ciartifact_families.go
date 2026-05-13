package ciartifact

func evaluateFamilies(reqs map[string]FamilyRequirement, inputs []FamilyInput) []FamilyObservation {
	// Family evaluation compares required and observed artifact groups without
	// collapsing individual family verdicts into an opaque aggregate.

	observed := observedFamilies(inputs)
	out, seen := requiredFamilyObservations(reqs, observed)
	for _, family := range extraFamilies(observed, seen) {
		out = append(out, evaluateFamily(FamilyRequirement{Family: family}, observed[family], false))
	}
	return out
}
