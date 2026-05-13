package managed

func missingManagedWitnessPassState(witness Witness) bool {
	// Witness status and freshness must both pass before bindings can be trusted.
	return witness.Status != StatePass || witness.FreshnessState != StatePass
}

func missingWitnessArtifacts(run RunEvidence, witness Witness) bool {
	// Artifact digests are the replayable link between selected output and witness.
	return len(run.OutputArtifacts) == 0 || len(witness.ArtifactDigests) == 0
}
