package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
)

func adapterCaptureExitCode(result adaptercapture.AssessmentResult) int {
	return stringExitCode(result.AdapterCaptureAssessment, adapterCaptureExitCodes, exitCannotVerify)
}
