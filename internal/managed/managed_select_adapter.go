package managed

func selectedAdapter(input Input) (Adapter, bool) {
	// selectedAdapter preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if input.Run.ManagedBoundaryEnrolled == nil {
		return Adapter{}, false
	}
	for _, adapter := range input.Registry.Adapters {
		if adapter.AdapterID == input.Run.ManagedBoundaryEnrolled.AdapterID {

			return adapter, true
		}
	}
	return Adapter{}, false
}

func selectedAuthorizedAdapter(input Input, adapter Adapter) (AuthorizedAdapter, bool) {
	// selectedAuthorizedAdapter preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, allowed := range input.Policy.AuthorizedAdapters {
		if authorizedAdapterMatches(allowed, adapter) {

			return allowed, true
		}
	}
	return AuthorizedAdapter{}, false
}

func authorizedAdapterMatches(allowed AuthorizedAdapter, adapter Adapter) bool {
	// authorizedAdapterMatches preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.

	return allowed.AdapterID == adapter.AdapterID &&
		allowed.HarnessID == adapter.HarnessID &&
		allowed.AuthorityRef == adapter.AuthorityRef &&
		allowed.DeploymentRef == adapter.DeploymentRef
}
