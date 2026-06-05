package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func reviewValidationExitCode(validation prreview.Validation) int {
	return stringExitCode(validation.ReviewCoverageState, reviewValidationExitCodes, 0)
}

var reviewValidationExitCodes = map[string]int{
	prreview.CoverageCannotVerify: exitCannotVerify,
	prreview.CoverageUnresolved:   exitCannotVerify,
}
