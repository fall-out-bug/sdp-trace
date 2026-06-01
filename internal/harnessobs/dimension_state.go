package harnessobs

func dimension(family string, required bool, count int) Dimension {
	// Observed families pass; absent families stay explicit, with required
	// absence represented as not_assessed rather than silently passing.
	if count > 0 {
		return Dimension{Family: family, Required: required, State: StatePass, ReasonCode: "event_family_observed", EventCount: count}
	}

	reason := "optional_event_family_absent"
	if required {
		reason = "required_event_family_absent"
	}
	return Dimension{Family: family, Required: required, State: StateNotAssessed, ReasonCode: reason}
}

func compose(dimensions []Dimension) (string, string) {
	// Composition considers required dimensions only; optional absence should
	// not degrade the validation state.
	// The highest-ranked required state supplies the validation reason so a
	// missing required family remains explainable in the final artifact.
	state := StatePass
	reason := "all_required_dimensions_observed"
	for _, dim := range dimensions {
		if !dim.Required {
			continue
		}

		if rank(dim.State) > rank(state) {
			state = dim.State
			reason = dim.ReasonCode
		}
	}
	return state, reason
}

func rank(state string) int {
	// Unknown states rank as zero through the map default, matching the legacy
	// behavior covered by harnessobs_crap_test.
	return stateRank[state]
}
