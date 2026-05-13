package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
)

func forensicExitCode(result forensic.AssessmentResult) int {
	return stringExitCode(result.ForensicRetentionAssessment, forensicExitCodes, exitCannotVerify)
}
