package harnessobs

import (
	"errors"
)

func validateSessionSetupAction(action SessionSetupAction) error {
	// validateSessionSetupAction keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !safeIDPattern.MatchString(action.ID) {
		return errors.New("unsafe setup action id")
	}

	switch action.Kind {
	case "init", "profile", "wrapper", "hook", "context_isolation":
		return nil
	default:
		return errors.New("unsupported setup action kind")
	}
}
