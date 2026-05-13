package adaptercapture

func modelIdentityCondition(run RunEvidence) Condition {
	// Model identity evidence preserves provider/model attribution before claims are
	// considered bound to the recorded work.
	for _, event := range run.AdapterEvents {
		if modelIdentityOverclaimed(run, event) {

			return fail("model_identity_not_overclaimed", "gateway_identity_overclaimed", "model identity is claimed as gateway-observed without bound gateway evidence", "Keep model identity harness_observed or bind gateway evidence.")
		}
	}
	if !run.GatewayIntegrated {

		return Condition{ID: "model_identity_not_overclaimed", State: StateNotIntegrated, ReasonCode: "gateway_not_integrated", Reason: "gateway provenance is not integrated", NextAction: "Integrate gateway evidence before claiming gateway-observed model identity."}
	}
	return pass("model_identity_not_overclaimed", "model_identity_not_overclaimed", "model identity provenance stays within available gateway evidence")
}

func modelIdentityOverclaimed(run RunEvidence, event AdapterEvent) bool {
	return event.EventType == "model_call_observed" && event.ModelIdentityProvenance == "gateway_observed" && !gatewayModelIdentityBound(run, event)
}

func gatewayModelIdentityBound(run RunEvidence, event AdapterEvent) bool {
	return run.GatewayIntegrated && run.GatewayEvidenceBound && event.IdentityBinding == IdentityBound
}
