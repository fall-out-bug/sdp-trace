package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func reportExitCode(summary demo.Summary) int {
	if summary.CannotVerifyCount > 0 {
		// Cannot-verify rows take precedence because missing evidence is not a
		// successful report even when no explicit failures were counted.
		return exitCannotVerify
	}
	if summary.FailedCount > 0 {
		return 1
	}
	return 0
}
