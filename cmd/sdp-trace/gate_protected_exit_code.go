package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

var protectedGateExitCodes = map[string]int{
	demo.GatePass:         0,
	demo.GateFail:         1,
	demo.GateCannotVerify: exitCannotVerify,
	demo.GateNotAssessed:  exitCannotVerify,
}

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
