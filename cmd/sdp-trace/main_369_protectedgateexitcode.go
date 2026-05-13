package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func protectedGateExitCode(result demo.GateResult) (int, bool) {
	if result.SelectedProfile != demo.GateProfileProtected {
		// Non-protected gates fall back to layered local/CI/audit state.
		return 0, false
	}
	code, ok := protectedGateExitCodes[result.ProtectedGate]
	if !ok {
		return 0, false
	}
	return code, true
}
