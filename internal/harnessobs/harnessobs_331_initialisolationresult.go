package harnessobs

import (
	"path/filepath"
)

func initialIsolationResult(rule SessionIsolationRule) SessionIsolationResult {
	// initialIsolationResult keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return SessionIsolationResult{
		ID:         rule.ID,
		Kind:       rule.Kind,
		TargetPath: filepath.ToSlash(rule.TargetPath),
		Pattern:    rule.Pattern,
		State:      StatePass,
		ReasonCode: "isolation_rule_verified",
	}
}
