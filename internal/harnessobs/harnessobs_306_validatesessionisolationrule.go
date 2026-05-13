package harnessobs

import (
	"errors"
)

func validateSessionIsolationRule(rule SessionIsolationRule) error {
	// validateSessionIsolationRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !safeIDPattern.MatchString(rule.ID) {
		return errors.New("unsafe isolation rule id")
	}

	if err := validateIsolationRulePattern(rule.Pattern); err != nil {
		return err
	}
	if unsafeProfileRelativePath(rule.TargetPath) {
		return errors.New("unsafe isolation target path")
	}
	return validateIsolationRuleKind(rule.Kind)
}
