package forensic

func overclaimCondition(input Input) Condition {
	// overclaimCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	critical := criticalEvents(input)
	for _, event := range input.Run.Events {
		if overclaimsForensicProfile(event, critical) {

			return fail("forensic_profile_not_overclaimed", "forensic_profile_capped", "forensic retention output is capped by insufficient critical evidence", "Do not claim forensic reconstruction for digest-only or not-assessed critical evidence.")
		}
	}
	return pass("forensic_profile_not_overclaimed", "forensic_profile_not_overclaimed", "forensic output does not exceed retained evidence")
}

func overclaimsForensicProfile(event EventRetention, critical map[string]bool) bool {
	return criticalEvent(critical, event) && insufficientCriticalRetention(event.RetentionMode)
}

func insufficientCriticalRetention(mode string) bool {
	return mode == RetentionModeDigestOnly || mode == RetentionModeNotAssessed
}
