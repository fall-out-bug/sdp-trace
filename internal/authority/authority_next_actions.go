package authority

func nextActions(result Result) []string {
	// nextActions keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if action, ok := authorityStateNextActions[result.AuthorityEvaluationState]; ok {
		return []string{action}
	}

	return []string{"Retain evidence references if downstream consumers need replay."}
}

var authorityStateNextActions = map[string]string{
	StateCannotVerify:     "Fix malformed, stale, inaccessible, or conflicting authority evidence before using these facts.",
	StateNotAssessed:      "Supply a selected policy_id, authority envelope, applicable rule, or required evidence before claiming authority compliance.",
	StateOutsideAuthority: "External policy consumers decide whether outside_authority blocks, contaminates, or requires escalation.",
}
