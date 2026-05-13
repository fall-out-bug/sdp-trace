package prreview

func attachPromptRef(result *ReviewerResult, role ReviewRole) error {
	// attachPromptRef keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	promptRef, err := promptSafeRef(role)
	if err != nil {

		result.Status = StatusCannotVerify
		result.Reason = "prompt_ref_cannot_verify"
		return err
	}
	result.PromptRef = promptRef
	return nil
}
