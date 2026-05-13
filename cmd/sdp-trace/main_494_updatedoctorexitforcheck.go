package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func updateDoctorExitForCheck(result string, exitCode int, check doctorCheck) (string, int) {
	if check.State == string(trace.VerdictCannotVerify) {
		// Any cannot-verify control point lowers the overall doctor exit.
		return string(trace.VerdictCannotVerify), exitCannotVerify
	}
	return result, exitCode
}
