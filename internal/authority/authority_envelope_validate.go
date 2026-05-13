package authority

func validateEnvelope(env AuthorityEnvelope) string {
	// validateEnvelope keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if env.PolicyID == "" || env.ActorRef == "" || env.TaskID == "" {

		return "authority_envelope_malformed"
	}
	return firstReason(
		validateEventSet(env.AllowedEvents, env.DeniedEvents),
		validateTargetRules(env),
		validateTargetRuleOverlap(env.TargetRules),
	)
}

func firstReason(reasons ...string) string {
	// firstReason keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	for _, reason := range reasons {
		if reason != "" {

			return reason
		}
	}
	return ""
}
