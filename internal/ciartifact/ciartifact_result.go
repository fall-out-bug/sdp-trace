package ciartifact

func observationResult(manifest Manifest, source SourceIdentity, run RunIdentity, inputs evaluatedInputs, state string, identityCannotVerify bool) ObservationResult {
	// Result assembly keeps assessment state, reasons, and next actions aligned.
	// The output is a replayable CI-artifact observation, not external CI proof.

	result := baseObservationResult(manifest, source, run, inputs, state)
	addObservationGaps(&result, inputs, identityCannotVerify)
	return result
}

func baseObservationResult(manifest Manifest, source SourceIdentity, run RunIdentity, inputs evaluatedInputs, state string) ObservationResult {
	// The base result begins in cannot_verify until concrete manifest evidence
	// raises or fails individual artifact-family conditions.
	result := ObservationResult{
		SchemaVersion:            SchemaVersion,
		SelectedProfile:          ProfileCIArtifactObservation,
		AuthorityScope:           safeAuthorityScope(manifest.AuthorityScope),
		ArtifactObservationState: state,

		SelectedSource: source,
		SelectedRun:    run,

		ProducerScope:       aggregateProducerScope(inputs.families),
		ArtifactAccessState: aggregateAccessState(inputs.families),
		RequiredFamilies:    orderedRequirements(inputs.reqs),
		ArtifactFamilies:    inputs.families,

		Bindings:      bindingSummary(inputs.families),
		ArtifactIndex: inputs.index,
		OutputSafety:  inputs.safety,

		SafetyRuleset: defaultSafetyRuleset(manifest.SafetyRuleset),
	}
	return result
}

func addObservationGaps(result *ObservationResult, inputs evaluatedInputs, identityCannotVerify bool) {
	// Gap rows are attached only after family and safety checks have run.
	// This avoids presenting skipped checks as successful artifact coverage.

	result.Reasons = reasons(inputs.families, inputs.index, inputs.safety, identityCannotVerify)
	result.NextActions = nextActions(inputs.families, inputs.index, inputs.safety, identityCannotVerify)
}
