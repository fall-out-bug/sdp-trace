package prreview

// prepareCommandRunner attaches prompt provenance and verifies a command exists
// before process execution.
func prepareCommandRunner(result *ReviewerResult, role ReviewRole) (*workingTreeBaseline, bool, error) {
	if err := attachPromptRef(result, role); err != nil {
		return nil, false, nil
	}

	return nil, commandConfigured(result, role), nil
}

// attachPromptRef records prompt provenance or marks the run cannot_verify when
// prompt references cannot be made portable.
func attachPromptRef(result *ReviewerResult, role ReviewRole) error {
	promptRef, err := promptSafeRef(role)
	if err != nil {
		result.Status = StatusCannotVerify
		result.Reason = "prompt_ref_cannot_verify"
		return err
	}
	result.PromptRef = promptRef
	return nil
}

// commandConfigured prevents an empty command from being treated as an executed
// review.
func commandConfigured(result *ReviewerResult, role ReviewRole) bool {
	if len(role.Command) == 0 {
		result.Reason = "runner_command_not_configured"
		return false
	}
	return true
}
