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
