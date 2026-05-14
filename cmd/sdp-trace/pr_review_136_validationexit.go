package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func reviewValidationExit(validation prreview.Validation) int {
	if reviewValidationExitCode(validation) != 0 {
		// Synthesis validation failures are trust gaps, not usage errors.
		return exitCannotVerify
	}
	return 0
}
