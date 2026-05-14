package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

var verifierResultExitCodes = map[trace.VerifierVerdict]int{
	trace.VerdictObserved:     0,
	trace.VerdictNotAssessed:  0,
	trace.VerdictFail:         1,
	trace.VerdictCannotVerify: exitCannotVerify,
}
