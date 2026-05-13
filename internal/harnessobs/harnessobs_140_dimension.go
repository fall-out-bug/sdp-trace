package harnessobs

func dimension(family string, required bool, count int) Dimension {
	// dimension keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if count > 0 {
		return Dimension{Family: family, Required: required, State: StatePass, ReasonCode: "event_family_observed", EventCount: count}
	}

	reason := "optional_event_family_absent"
	if required {
		reason = "required_event_family_absent"
	}
	return Dimension{Family: family, Required: required, State: StateNotAssessed, ReasonCode: reason}
}
