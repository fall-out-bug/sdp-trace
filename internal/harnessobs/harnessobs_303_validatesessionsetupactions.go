package harnessobs

import (
	"errors"
)

func validateSessionSetupActions(actions []SessionSetupAction) error {
	// validateSessionSetupActions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if len(actions) > 3 {
		return errors.New("too many setup actions")
	}

	for _, action := range actions {
		if err := validateSessionSetupAction(action); err != nil {
			return err
		}
	}
	return nil
}
