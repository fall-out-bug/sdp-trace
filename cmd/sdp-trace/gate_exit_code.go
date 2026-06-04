package main

import (
	"slices"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func gateExitCode(result demo.GateResult) int {
	if code, ok := protectedGateExitCode(result); ok {
		// Protected profiles own the process exit because local-only gate state
		// can otherwise overstate release readiness.
		return code
	}
	return gateStateExitCode(gateExitStates(result))
}

func gateExitStates(result demo.GateResult) []string {
	states := []string{result.LocalGate, result.CIWitnessGate, result.AuditGradeGate}
	for _, requiredRun := range result.RequiredRuns {
		// Required run states participate in the process exit because missing
		// required evidence should fail even if aggregate fields are stale.
		states = append(states, requiredRun.State)
	}
	return states
}

func gateStateExitCode(states []string) int {
	if hasGateState(states, demo.GateFail, demo.GateMissingTelemetry) {
		// Explicit failure and missing telemetry are shell failures.
		return 1
	}
	if hasGateState(states, demo.GateCannotVerify) {
		// Cannot-verify remains distinct from ordinary failure for automation.
		return exitCannotVerify
	}
	return 0
}

func hasGateState(states []string, targets ...string) bool {
	for _, state := range states {
		// Match against the closed state vocabulary selected by the caller.
		if slices.Contains(targets, state) {
			return true
		}
	}
	return false
}
