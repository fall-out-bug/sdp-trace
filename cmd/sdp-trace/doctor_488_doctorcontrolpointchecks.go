package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

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
