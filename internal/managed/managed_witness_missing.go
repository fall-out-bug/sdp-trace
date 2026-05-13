package managed

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
