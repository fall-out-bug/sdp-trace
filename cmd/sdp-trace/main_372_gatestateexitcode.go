package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

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
