package adaptercapture

func providerRefsCondition(run RunEvidence) Condition {
	// Provider refs are checked for unsafe material before any external reference is
	// echoed into evidence output.
	if providerRefsContainSecret(run.ProviderRefs) {

		return fail("provider_refs_portable", "provider_ref_contains_secret", "provider-neutral reference contains credential-like material", "Persist canonical token-free provider references.")
	}
	if adapterEventsProviderRefsContainSecret(run.AdapterEvents) {

		return fail("provider_refs_portable", "provider_ref_contains_secret", "event-level provider reference contains credential-like material", "Persist canonical token-free provider references.")
	}
	return pass("provider_refs_portable", "provider_refs_portable", "provider references are portable and token-free")
}

func providerRefsContainSecret(refs []ProviderRef) bool {
	// providerRefsContainSecret preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	for _, ref := range refs {
		if providerRefContainsSecret(ref) {

			return true
		}
	}
	return false
}

func adapterEventsProviderRefsContainSecret(events []AdapterEvent) bool {
	// adapterEventsProviderRefsContainSecret preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	for _, event := range events {
		if eventProviderRefsContainSecret(event) {

			return true
		}
	}
	return false
}

func providerRefContainsSecret(ref ProviderRef) bool {
	return containsSecret(ref.SourceRef) || containsSecret(ref.ChangeRef) || containsSecret(ref.ReviewRef)
}

func eventProviderRefsContainSecret(event AdapterEvent) bool {
	return stringSliceContainsSecret(event.ProviderRefs)
}
