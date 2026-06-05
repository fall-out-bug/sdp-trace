package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func doctorEnvironmentChecks() []doctorCheck {
	// Environment checks describe what the local process can know; they are not
	// external witness evidence.
	return []doctorCheck{
		{
			ID:     "local_process",
			State:  "pass",
			Reason: "current process can inspect local environment",
		},
		{
			ID:     "offline_development",
			State:  "offline_dev",
			Reason: "external CI identity is not required for local preview or wrapper readiness",
		},
	}
}

func doctorControlPointChecks(defaultContract trace.Contract, ciCheck doctorCheck, checks ...doctorCheck) []doctorCheck {
	// Control-point order is stable because downstream fixtures and reports cite
	// named gaps in the order doctor prints them.
	controlPoints := []doctorCheck{
		{
			ID:     "local_wrapper",
			State:  "pass",
			Reason: "wrap and run commands are registered in this binary",
		},
	}
	controlPoints = append(controlPoints, checks...)
	// The built-in contract is reported before CI prerequisites so local
	// readiness remains separate from unavailable external identity.
	controlPoints = append(controlPoints, doctorDefaultContractCheck(defaultContract))
	return append(controlPoints, ciCheck)
}

func updateDoctorExitForCheck(result string, exitCode int, check doctorCheck) (string, int) {
	if check.State == string(trace.VerdictCannotVerify) {
		// Any cannot-verify control point lowers the overall doctor exit.
		return string(trace.VerdictCannotVerify), exitCannotVerify
	}
	return result, exitCode
}
