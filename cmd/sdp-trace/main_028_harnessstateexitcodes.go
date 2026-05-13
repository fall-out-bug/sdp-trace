package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

var harnessStateExitCodes = map[string]int{
	harnessobs.StatePass:         0,
	harnessobs.StateFail:         1,
	harnessobs.StateNotAssessed:  exitCannotVerify,
	harnessobs.StateCannotVerify: exitCannotVerify,
}
