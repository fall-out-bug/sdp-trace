package verifier

import "github.com/fall_out_bug/sdp-trace/internal/trace"

func observedResult(runDir, runID string) trace.VerifierResult {
	// Observed is the optimistic local replay baseline before requirement gaps
	// are applied.
	return trace.VerifierResult{
		RunID:         runID,
		RunDir:        runDir,
		TrustScope:    trace.TrustScopeLocalObserved,
		Result:        trace.VerdictObserved,
		Completeness:  trace.CompletenessComplete,
		Replayability: trace.ReplayabilityPartial,
	}
}

func cannotVerifyResult(runDir, runID, reason string) trace.VerifierResult {
	// Cannot-verify results still carry local trust scope because the verifier
	// only assessed local artifacts.
	return trace.VerifierResult{
		RunID:        runID,
		RunDir:       runDir,
		Result:       trace.VerdictCannotVerify,
		TrustScope:   trace.TrustScopeLocalObserved,
		Completeness: trace.CompletenessUnknown,
		Reason:       reason,
	}
}

func cannotVerifyReplayResult(result trace.VerifierResult, reason string) trace.VerifierResult {
	// Preserve run identity while lowering replayability and completeness.
	result.Result = trace.VerdictCannotVerify
	result.Completeness = trace.CompletenessUnknown
	result.Replayability = trace.ReplayabilityNone
	result.Reason = reason
	return result
}

func resultWithMissingEvidence(result trace.VerifierResult, contract trace.Contract, events []trace.Event) (trace.VerifierResult, trace.MissingEvidenceTable) {
	missingEvidence := trace.GenerateMissingEvidenceTable(contract, observedEventSet(events))
	if len(missingEvidence.Rows) > 0 {
		// Missing required evidence is not_assessed: the replay succeeded, but
		// the contract is not fully covered.
		result.Result = trace.VerdictNotAssessed
		result.Completeness = trace.CompletenessMissingTelemetry
		result.Replayability = trace.ReplayabilityPartial
	}
	return result, missingEvidence
}
