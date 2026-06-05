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

func verifierResultExitCode(result trace.VerifierVerdict) int {
	code, ok := verifierResultExitCodes[result]
	if !ok {
		// Unknown future verifier verdicts should not fail old automation by
		// default; schema validation remains the stronger compatibility check.
		return 0
	}
	return code
}
