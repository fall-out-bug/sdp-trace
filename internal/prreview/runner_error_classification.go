package prreview

// applyRunnerError classifies process and prompt failures into explicit
// reviewer evidence states.
func applyRunnerError(result *ReviewerResult, err error) error {
	if err == nil {
		return nil
	}
	status, reason := runnerErrorState(err)
	result.Status = status
	result.Reason = reason
	return err
}

// runnerErrorState keeps prompt replay failures separate from missing runner
// configuration and failed runner execution.
func runnerErrorState(err error) (string, string) {
	if promptTemplateCannotVerify(err) {
		return StatusCannotVerify, "prompt_ref_cannot_verify"
	}
	if promptEvidenceCannotVerify(err) {
		return StatusCannotVerify, "prompt_evidence_cannot_verify"
	}
	if runnerUnavailable(err) {
		return StatusNotAssessed, "runner_unavailable"
	}
	return StatusFailed, "runner_failed"
}
