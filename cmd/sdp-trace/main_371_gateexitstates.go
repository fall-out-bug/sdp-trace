package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func gateExitStates(result demo.GateResult) []string {
	states := []string{result.LocalGate, result.CIWitnessGate, result.AuditGradeGate}
	for _, requiredRun := range result.RequiredRuns {
		// Required run states participate in the process exit because missing
		// required evidence should fail even if aggregate fields are stale.
		states = append(states, requiredRun.State)
	}
	return states
}
