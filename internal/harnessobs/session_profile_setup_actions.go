package harnessobs

import "errors"

func validateSessionSetupActions(actions []SessionSetupAction) error {
	// Setup actions are intentionally capped so a profile cannot expand into
	// an unbounded local setup workflow before observation starts.
	if len(actions) > 3 {
		return errors.New("too many setup actions")
	}

	for _, action := range actions {
		// Preserve first-error ordering for profile authors.
		if err := validateSessionSetupAction(action); err != nil {
			return err
		}
	}
	return nil
}

func validateSessionSetupAction(action SessionSetupAction) error {
	// Action IDs are later copied into trace output, so they must be safe
	// identifiers before kind-specific behavior is considered.
	if !safeIDPattern.MatchString(action.ID) {
		return errors.New("unsafe setup action id")
	}

	switch action.Kind {
	case "init", "profile", "wrapper", "hook", "context_isolation":
		// These names describe declared setup intent only; execution remains in
		// the dedicated session setup path.
		return nil
	default:
		return errors.New("unsupported setup action kind")
	}
}
