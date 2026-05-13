package managed

func adapterActivationCondition(input Input) Condition {
	// adapterActivationCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if input.Run.ManagedBoundaryEnrolled == nil || input.Run.ManagedBoundaryEnrolled.AdapterID == "" {

		return cannotVerify("adapter_activation_observed", "adapter_activation_missing", "adapter activation cannot be verified", "Record adapter activation before child launch.")
	}
	return pass("adapter_activation_observed", "adapter_activation_observed", "adapter activation is bound to managed enrollment")
}

func adapterConnectionCondition(input Input) Condition {
	// adapterConnectionCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if input.Run.AdapterDisconnectObserved {
		return fail("adapter_connection_continuous", "adapter_disconnect_observed", "adapter disconnected during required managed observation window", "Rerun with continuous managed adapter connection.")
	}

	return pass("adapter_connection_continuous", "adapter_connection_continuous", "adapter connection has no observed disconnect during required window")
}
