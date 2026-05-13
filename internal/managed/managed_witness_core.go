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
