package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
	"github.com/fall_out_bug/sdp-trace/internal/managed"
)

func adapterCaptureExitCode(result adaptercapture.AssessmentResult) int {
	return stringExitCode(result.AdapterCaptureAssessment, adapterCaptureExitCodes, exitCannotVerify)
}

func managedExitCode(result managed.AssessmentResult) int {
	return stringExitCode(result.ManagedHarnessAssessment, managedExitCodes, exitCannotVerify)
}

func forensicExitCode(result forensic.AssessmentResult) int {
	return stringExitCode(result.ForensicRetentionAssessment, forensicExitCodes, exitCannotVerify)
}
