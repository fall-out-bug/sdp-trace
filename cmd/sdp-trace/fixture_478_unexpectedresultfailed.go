package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func unexpectedFixtureResultFailed(result trace.VerifierResult) bool {
	// Without an explicit expected result, only fail/cannot_verify are treated as
	// fixture failures; observed/not_assessed remain inspectable but nonfatal.
	return result.Result == trace.VerdictFail || result.Result == trace.VerdictCannotVerify
}
