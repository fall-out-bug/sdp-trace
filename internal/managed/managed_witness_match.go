package managed

func managedWitnessBindingMismatch(input Input, boundary ManagedBoundaryEnrolled) bool {
	// managedWitnessBindingMismatch preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.

	return !witnessMatchesRun(input.Witness, input.Run) ||
		!witnessMatchesAuthority(input.Witness, input.Policy, input.Registry) ||
		!witnessMatchesEvents(input.Witness, input.Run, boundary) ||
		!artifactsMatch(input.Run.OutputArtifacts, input.Witness.ArtifactDigests)
}

func witnessMatchesRun(witness Witness, run RunEvidence) bool {
	return witnessRunIdentityMatches(witness, run) &&
		witnessRunTraceMatches(witness, run)
}

func witnessRunIdentityMatches(witness Witness, run RunEvidence) bool {
	return witness.RunID == run.RunID && witness.RunNonce == run.RunNonce && witness.SourceCommit == run.SourceCommit
}

func witnessRunTraceMatches(witness Witness, run RunEvidence) bool {
	return witness.ChainHead == run.ChainHead &&
		witness.EventCount == run.EventCount
}

func witnessMatchesAuthority(witness Witness, policy Policy, registry Registry) bool {
	return witness.ManagedPolicyDigest == policy.PolicyProvenance.Digest &&
		witness.AdapterRegistryDigest == registry.Provenance.Digest
}

func witnessMatchesEvents(witness Witness, run RunEvidence, boundary ManagedBoundaryEnrolled) bool {
	return witness.EnrollmentEventDigest == boundary.EventDigest &&
		witness.LaunchEventDigest == run.ChildLaunch.EventDigest
}
