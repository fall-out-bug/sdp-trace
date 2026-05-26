package main

import (
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
