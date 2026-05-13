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
	if errors.Is(err, exec.ErrNotFound) || strings.Contains(err.Error(), "executable file not found") {

		result.Status = StatusNotAssessed
		result.Reason = "runner_unavailable"
	} else {

		result.Status = StatusFailed
		result.Reason = "runner_failed"
	}
	return err
}
