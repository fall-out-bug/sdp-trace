package adaptercapture

func identityCondition(run RunEvidence) Condition {
	// Identity evidence binds provider, model, and adapter labels before event claims
	// can contribute to a capture verdict.
	for _, event := range run.AdapterEvents {
		if condition, ok := identityConditionForEvent(event); ok {

			return condition
		}
	}
	return pass("adapter_identity_visible", "adapter_identity_visible", "adapter and producer identity are visible with binding classification")
}

func identityConditionForEvent(event AdapterEvent) (Condition, bool) {
	// identityConditionForEvent preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	if adapterIdentityMissing(event) {

		return cannotVerify("adapter_identity_visible", "adapter_identity_missing", "adapter or producer identity is missing", "Record adapter and producer identity."), true
	}
	if !validIdentityBinding(event.IdentityBinding) {

		return cannotVerify("adapter_identity_visible", "adapter_identity_unclassified", "adapter identity binding state is not classified", "Classify adapter identity as self_asserted or bound."), true
	}
	return Condition{}, false
}

func adapterIdentityMissing(event AdapterEvent) bool {
	return event.ProducerIdentity == "" || event.AdapterIdentity == ""
}

func validIdentityBinding(binding string) bool {
	return binding == IdentitySelfAsserted || binding == IdentityBound
}
