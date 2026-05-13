package verifier

import "github.com/fall_out_bug/sdp-trace/internal/trace"

func failedChainResult(runDir, runID, issue string) trace.VerifierResult {
	// A broken event chain is fail evidence because replay contradicted retained
	// source artifacts. It is not downgraded to cannot_verify.
	return trace.VerifierResult{
		RunID:         runID,
		RunDir:        runDir,
		TrustScope:    trace.TrustScopeLocalObserved,
		Result:        trace.VerdictFail,
		Completeness:  trace.CompletenessMissingTelemetry,
		Replayability: trace.ReplayabilityNone,
		Reason:        issue,
	}
}
