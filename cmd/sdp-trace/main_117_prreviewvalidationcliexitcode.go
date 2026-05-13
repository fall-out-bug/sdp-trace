package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func prReviewValidationCLIExitCode(validation prreview.Validation) int {
	if reviewValidationExitCode(validation) != 0 {
		// Invalid review evidence cannot support a PR trust claim.
		return exitCannotVerify
	}
	// A zero exit here only means the review packet validated locally.
	return 0
}
