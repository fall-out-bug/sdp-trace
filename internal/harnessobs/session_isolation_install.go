package harnessobs

import "errors"

// installIsolationRules resolves profile-relative targets before mutating files,
// so the returned verification evidence always points at the concrete target.
func installIsolationRules(profilePath string, rules []SessionIsolationRule) ([]SessionIsolationResult, error) {
	results := make([]SessionIsolationResult, 0, len(rules))
	for _, rule := range rules {
		result, err := installProfileIsolationRule(profilePath, rule)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// installProfileIsolationRule keeps target resolution and materialization in
// one per-rule path while the caller handles aggregation.
func installProfileIsolationRule(profilePath string, rule SessionIsolationRule) (SessionIsolationResult, error) {
	resolvedRule, err := resolveIsolationRuleTarget(profilePath, rule)
	if err != nil {
		return SessionIsolationResult{}, err
	}

	return installIsolationRule(resolvedRule)
}

// resolveIsolationRuleTarget keeps profile-local rule paths relative to the
// session profile while rejecting traversal and unsafe filenames.
func resolveIsolationRuleTarget(profilePath string, rule SessionIsolationRule) (SessionIsolationRule, error) {
	targetPath, err := safeProfileRelativeIsolationFile(profilePath, rule.TargetPath)
	if err != nil {
		return SessionIsolationRule{}, err
	}
	rule.TargetPath = targetPath
	return rule, nil
}

// installIsolationRule materializes one rule and immediately verifies the
// readback evidence used in the reported isolation result.
func installIsolationRule(rule SessionIsolationRule) (SessionIsolationResult, error) {
	if err := ensureIsolationRule(rule); err != nil {
		return SessionIsolationResult{}, err
	}

	return verifyIsolationRule(rule)
}

// ensureIsolationRule dispatches only supported rule kinds; unknown kinds stay
// errors rather than being reported as unverifiable evidence.
func ensureIsolationRule(rule SessionIsolationRule) error {
	installer, ok := isolationRuleInstallers[rule.Kind]
	if !ok {
		return errors.New("unsupported isolation rule kind")
	}
	return installer(rule.TargetPath, rule.Pattern)
}
