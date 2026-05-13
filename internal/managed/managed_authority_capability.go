package managed

func capabilityCondition(input Input) Condition {
	// Capability checks compare policy requirements to selected adapter claims
	// without assuming adapter self-description is sufficient proof.
	adapter, ok := selectedAdapter(input)
	if !ok {

		return cannotVerify("adapter_capabilities_satisfy_contract", "adapter_capability_missing", "selected adapter capabilities cannot be verified", "Supply an authorized adapter with declared capabilities.")
	}
	return selectedAdapterCapabilityCondition(input, adapter)
}

func selectedAdapterCapabilityCondition(input Input, adapter Adapter) Condition {
	// selectedAdapterCapabilityCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	authorized, condition, ok := managedCapabilityPolicy(input, adapter)
	if !ok {
		return condition
	}
	if !adapterSatisfiesPolicyCapabilities(adapter, authorized) {

		return cannotVerify("adapter_capabilities_satisfy_contract", "adapter_capability_missing", "adapter capability references do not satisfy managed policy", "Use an adapter whose declared capabilities match managed policy requirements.")
	}
	if !adapterCapabilitiesCoverEvents(input, adapter, authorized) {

		return cannotVerify("adapter_capabilities_satisfy_contract", "adapter_capability_missing", "adapter capability set does not cover a required event type", "Use an adapter whose capabilities cover the managed contract.")
	}
	return pass("adapter_capabilities_satisfy_contract", "adapter_capabilities_satisfy_contract", "adapter capabilities cover required event types")
}
func managedCapabilityPolicy(input Input, adapter Adapter) (AuthorizedAdapter, Condition, bool) {
	// managedCapabilityPolicy preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	authorized, ok := selectedAuthorizedAdapter(input, adapter)
	if !ok || len(authorized.CapabilityIDs) == 0 {

		return AuthorizedAdapter{}, cannotVerify("adapter_capabilities_satisfy_contract", "adapter_capability_missing", "managed policy does not name required adapter capabilities", "Supply managed policy capability requirements for the selected adapter."), false
	}
	return authorized, Condition{}, true
}
