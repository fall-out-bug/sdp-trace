package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/managed"
)

func managedExitCode(result managed.AssessmentResult) int {
	return stringExitCode(result.ManagedHarnessAssessment, managedExitCodes, exitCannotVerify)
}
