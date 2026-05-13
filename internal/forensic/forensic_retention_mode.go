package forensic

func retentionModeCondition(input Input) Condition {
	// retentionModeCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	allowed := allowedRetentionModes(input.Policy)
	for _, event := range input.Run.Events {

		if condition, ok := retentionModeConditionForEvent(event, allowed); ok {
			return condition
		}
	}
	return pass("retention_mode_declared", "retention_mode_declared", "events declare FR-054 retention modes")
}

func retentionModeConditionForEvent(event EventRetention, allowed map[string]bool) (Condition, bool) {
	// retentionModeConditionForEvent keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if !validRetentionMode(event.RetentionMode) {

		return fail("retention_mode_declared", "invalid_retention_mode", "event declares a non-FR-054 retention mode", "Use digest_only, sanitized_excerpt, encrypted_raw_ref, external_artifact_ref, or not_assessed."), true
	}
	if len(allowed) > 0 && !allowed[event.RetentionMode] {

		return fail("retention_mode_declared", "retention_mode_not_policy_allowed", "event retention mode is not allowed by the selected redaction policy", "Use a retention mode allowed by the selected policy."), true
	}
	return Condition{}, false
}
