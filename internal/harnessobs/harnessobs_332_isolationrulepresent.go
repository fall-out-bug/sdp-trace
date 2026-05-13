package harnessobs

import (
	"errors"
)

func isolationRulePresent(rule SessionIsolationRule) (bool, error) {
	// isolationRulePresent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch rule.Kind {
	case "ignore_line":
		return lineIsolationRulePresent(rule.TargetPath, rule.Pattern)
	case "json_read_deny":
		return jsonReadDenyRulePresent(rule.TargetPath, rule.Pattern)
	default:

		return false, errors.New("unsupported isolation rule kind")
	}
}
