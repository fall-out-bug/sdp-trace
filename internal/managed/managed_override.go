package managed

func overrideCondition(input Input) Condition {
	// Override evidence is terminal closure context; it does not erase earlier
	// adapter, event, or witness failures.
	if input.Run.OverrideAttemptsTrustUpgrade {

		return fail("override_does_not_upgrade_managed_profile", "override_upgrade_rejected", "override artifact attempts to upgrade managed profile state", "Record override as non-upgrading evidence only.")
	}
	if input.Run.OverridePresent {
		return pass("override_does_not_upgrade_managed_profile", "override_present_non_upgrading", "override request is visible and non-upgrading")
	}
	return pass("override_does_not_upgrade_managed_profile", "override_absent_non_upgrading", "no override request is available to upgrade managed profile")
}
