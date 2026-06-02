package harnessobs

import (
	"errors"
	"strings"
)

func validateSessionIsolationRules(rules []SessionIsolationRule) error {
	for _, rule := range rules {
		// Keep validation deterministic by returning the first unsafe rule.
		if err := validateSessionIsolationRule(rule); err != nil {
			return err
		}
	}
	return nil
}

func validateSessionIsolationRule(rule SessionIsolationRule) error {
	// Rule IDs are emitted in setup results, so validate identity before
	// interpreting patterns or target paths.
	if !safeIDPattern.MatchString(rule.ID) {
		return errors.New("unsafe isolation rule id")
	}

	if err := validateIsolationRulePattern(rule.Pattern); err != nil {
		return err
	}
	// Target paths are only profile-relative declarations here; filesystem
	// resolution and installation are owned by the next slice.
	if unsafeProfileRelativePath(rule.TargetPath) {
		return errors.New("unsafe isolation target path")
	}
	return validateIsolationRuleKind(rule.Kind)
}

func validateIsolationRulePattern(pattern string) error {
	// Empty or multi-line patterns are rejected before they can affect line or
	// JSON read-deny materialization.
	if unsafeIsolationRulePattern(pattern) {
		return errors.New("unsafe isolation rule pattern")
	}
	return nil
}

func unsafeIsolationRulePattern(pattern string) bool {
	return strings.TrimSpace(pattern) == "" || strings.Contains(pattern, "\n") || strings.Contains(pattern, "\r")
}

func validateIsolationRuleKind(kind string) error {
	switch kind {
	case "ignore_line", "json_read_deny":
		// Supported kinds map to package-local installers; unknown kinds must
		// not silently become no-ops.
		return nil
	default:
		return errors.New("unsupported isolation rule kind")
	}
}
