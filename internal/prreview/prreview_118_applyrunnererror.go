package prreview

import (
	"errors"

	"os/exec"

	"strings"
)

func applyRunnerError(result *ReviewerResult, err error) error {
	// applyRunnerError keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if err == nil {
		return nil
	}
	status, reason := runnerErrorState(err)
	result.Status = status
	result.Reason = reason
	return err
}

func runnerErrorState(err error) (string, string) {
	// Prompt evidence failures mean the review packet cannot be replayed.
	// Missing runner binaries are configuration gaps, so they stay not_assessed.
	// Other process failures are failed reviewer executions.
	if errors.Is(err, errPromptTemplateCannotVerify) {
		return StatusCannotVerify, "prompt_ref_cannot_verify"
	}
	if errors.Is(err, errPromptEvidenceCannotVerify) {
		return StatusCannotVerify, "prompt_evidence_cannot_verify"
	}
	if errors.Is(err, exec.ErrNotFound) || strings.Contains(err.Error(), "executable file not found") {
		return StatusNotAssessed, "runner_unavailable"
	}
	return StatusFailed, "runner_failed"
}
