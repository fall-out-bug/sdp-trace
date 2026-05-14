package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
)

func checkpointVerifyExitCode(state string) int {
	if state == checkpoint.StatePass {
		// Only an explicit checkpoint pass maps to shell success.
		return 0
	}
	if state == checkpoint.StateCannotVerify {
		return exitCannotVerify
	}
	return 1
}
