package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func verifierResultExitCode(result trace.VerifierVerdict) int {
	code, ok := verifierResultExitCodes[result]
	if !ok {
		// Unknown future verifier verdicts should not fail old automation by
		// default; schema validation remains the stronger compatibility check.
		return 0
	}
	return code
}
