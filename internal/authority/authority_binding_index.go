package authority

func bindingStatesByEvent(bindings []EvidenceBinding) map[string][]EvidenceBinding {
	// bindingStatesByEvent keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	out := map[string][]EvidenceBinding{}
	for _, binding := range bindings {

		out[binding.LeftEventID] = append(out[binding.LeftEventID], binding)
		out[binding.RightEventID] = append(out[binding.RightEventID], binding)
	}
	return out
}

func hasVerifiedGatewayBinding(action ObservedAction, bindings []EvidenceBinding) bool {
	// hasVerifiedGatewayBinding keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	for _, binding := range bindings {
		if binding.BindingState == BindingVerified && binding.BindingType == "same_gateway_request" {

			return true
		}
	}
	return false
}

func bindingCannotVerify(bindings []EvidenceBinding) bool {
	// bindingCannotVerify keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	for _, binding := range bindings {
		if binding.BindingState == BindingCannotVerify {

			return true
		}
	}
	return false
}
