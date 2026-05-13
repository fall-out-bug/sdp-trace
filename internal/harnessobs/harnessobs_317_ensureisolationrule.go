package harnessobs

import (
	"errors"
)

func ensureIsolationRule(rule SessionIsolationRule) error {
	// ensureIsolationRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	installer, ok := isolationRuleInstallers[rule.Kind]
	if !ok {

		return errors.New("unsupported isolation rule kind")
	}
	return installer(rule.TargetPath, rule.Pattern)
}
