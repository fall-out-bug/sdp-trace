package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/authority"
)

func authorityExitCode(result authority.Result) int {
	return stringExitCode(result.AuthorityEvaluationState, authorityExitCodes, exitCannotVerify)
}
