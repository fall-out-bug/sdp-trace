package managed

func witnessCondition(input Input) Condition {
	// Witness evidence binds managed output back to run/report artifacts rather
	// than trusting a checked-in witness record by itself.
	witness := input.Witness
	if condition, ok := missingManagedWitnessCondition(input.Run, witness); ok {
		return condition
	}
	if managedWitnessMismatches(input) {

		return fail("managed_witness_bound", "managed_witness_mismatch", "managed witness does not bind the selected run, policy, registry, chain, or artifacts", "Regenerate managed witness evidence for the selected run.")
	}
	return pass("managed_witness_bound", "managed_witness_bound", "managed witness binds source, run, policy, registry, chain, and artifacts")
}

func missingManagedWitnessCondition(run RunEvidence, witness Witness) (Condition, bool) {
	// missingManagedWitnessCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if witness.WitnessID == "" {

		return cannotVerify("managed_witness_bound", "managed_witness_missing", "managed witness evidence is required", "Supply managed witness evidence bound to the run."), true
	}
	return invalidManagedWitnessCondition(run, witness)
}

func invalidManagedWitnessCondition(run RunEvidence, witness Witness) (Condition, bool) {
	// invalidManagedWitnessCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	switch {
	case missingManagedWitnessPassState(witness):

		return cannotVerify("managed_witness_bound", "managed_witness_missing", "managed witness is missing pass/freshness state", "Supply fresh managed witness evidence."), true
	case missingWitnessArtifacts(run, witness):

		return cannotVerify("managed_witness_bound", "managed_witness_missing", "managed witness artifact binding is required", "Supply managed witness evidence with output artifact digests."), true
	default:
		return Condition{}, false
	}
}

func missingManagedWitnessPassState(witness Witness) bool {
	return witness.Status != StatePass || witness.FreshnessState != StatePass
}

func missingWitnessArtifacts(run RunEvidence, witness Witness) bool {
	return len(run.OutputArtifacts) == 0 || len(witness.ArtifactDigests) == 0
}

func managedWitnessMismatches(input Input) bool {
	// managedWitnessMismatches preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	boundary := input.Run.ManagedBoundaryEnrolled
	if boundary == nil {

		return true
	}
	return managedWitnessBindingMismatch(input, *boundary)
}
