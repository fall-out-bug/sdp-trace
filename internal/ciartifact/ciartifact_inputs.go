package ciartifact

type evaluatedInputs struct {
	reqs     map[string]FamilyRequirement
	families []FamilyObservation
	index    ArtifactIndexResult
	safety   OutputSafetyResult
}

func evaluatedManifestInputs(manifest Manifest) evaluatedInputs {
	// Manifest input evaluation is the boundary where raw manifest fields become
	// normalized artifact-family observations.

	reqs := requirements(manifest.RequiredFamilies)
	return evaluatedInputs{
		reqs:     reqs,
		families: evaluateFamilies(reqs, manifest.ArtifactFamilies),
		index:    evaluateIndex(manifest.ArtifactIndex),
		safety:   evaluateSafety(manifest.OutputSafety),
	}
}

func requirements(input []FamilyRequirement) map[string]FamilyRequirement {
	// Requirements are derived from declared producer and access expectations,
	// keeping absent declarations distinct from failed declarations.
	reqs := map[string]FamilyRequirement{}
	for _, req := range input {

		family := canonicalFamily(req.Family)
		if family == "" {
			continue
		}
		req.Family = family
		req.RequiredProducerScope = safeRequiredProducerScope(req.RequiredProducerScope)
		reqs[family] = req
	}
	return reqs
}
