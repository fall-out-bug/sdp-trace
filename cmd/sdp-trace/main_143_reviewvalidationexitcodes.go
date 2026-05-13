package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

var reviewValidationExitCodes = map[string]int{
	prreview.CoverageCannotVerify: exitCannotVerify,
	prreview.CoverageUnresolved:   exitCannotVerify,
}
