package managed

func adapterIdentityCondition(input Input) Condition {
	// adapterIdentityCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	adapter, ok := selectedAdapter(input)
	if !ok {

		return fail("adapter_identity_authorized", "adapter_identity_unauthorized", "selected adapter is not present in the registry", "Register an authorized adapter before the run.")
	}
	if !adapterIdentityAuthorized(input, adapter) {
		return fail("adapter_identity_authorized", "adapter_identity_unauthorized", "adapter identity is not verified and authorized by policy", "Use a verified adapter identity authorized by managed policy.")
	}
	return pass("adapter_identity_authorized", "adapter_identity_authorized", "adapter identity is verified and authorized by policy")
}

func adapterIdentityAuthorized(input Input, adapter Adapter) bool {
	// adapterIdentityAuthorized preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if adapter.IdentityState != IdentityVerified {

		return false
	}
	_, ok := selectedAuthorizedAdapter(input, adapter)
	return ok
}
