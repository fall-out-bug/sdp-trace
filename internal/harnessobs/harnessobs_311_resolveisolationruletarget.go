package harnessobs

func resolveIsolationRuleTarget(profilePath string, rule SessionIsolationRule) (SessionIsolationRule, error) {
	// resolveIsolationRuleTarget keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	targetPath, err := safeProfileRelativeIsolationFile(profilePath, rule.TargetPath)
	if err != nil {
		return SessionIsolationRule{}, err
	}
	rule.TargetPath = targetPath
	return rule, nil
}
