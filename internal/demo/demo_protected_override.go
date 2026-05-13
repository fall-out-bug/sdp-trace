package demo

func protectedOverrideCondition(overrides []OverrideRequest) ProtectedCondition {
	// protectedOverrideCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if len(overrides) == 0 {
		return ProtectedCondition{ID: "override_does_not_upgrade_profile", State: GatePass, ReasonCode: "no_override_present", Reason: "no override request is available to upgrade the profile"}
	}
	for _, override := range overrides {

		if override.State != GatePass {

			return ProtectedCondition{
				ID:         "override_does_not_upgrade_profile",
				State:      GateCannotVerify,
				ReasonCode: "override_cannot_verify_non_upgrading",
				Reason:     "override request cannot verify and remains non-upgrading",
				NextAction: "Inspect override request evidence outside protected gate evaluation.",
			}
		}
	}
	return ProtectedCondition{ID: "override_does_not_upgrade_profile", State: GatePass, ReasonCode: "override_visible_non_upgrading", Reason: "override request is visible and non-upgrading"}
}
