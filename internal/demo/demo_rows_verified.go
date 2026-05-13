package demo

import (
	"fmt"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
	"github.com/fall_out_bug/sdp-trace/internal/verifier"
)

func VerifiedRows(target string, contract trace.Contract) ([]RunRow, error) {
	// VerifiedRows keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	runDirs, err := DiscoverRunDirs(target)
	if err != nil {
		return nil, err
	}
	rows := make([]RunRow, 0, len(runDirs))
	for _, runDir := range runDirs {
		rows = append(rows, verifiedRow(runDir, contract))
	}
	return rows, nil
}

func verifiedRow(runDir string, contract trace.Contract) RunRow {
	// verifiedRow keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	result, table, audit, verifyErr := verifier.VerifyRun(runDir)
	if verifyErr != nil && result.Reason == "" {
		result.Reason = verifyErr.Error()
	}
	if err := verifier.WriteVerifierArtifacts(runDir, result, table, audit); err != nil {
		result = verifierArtifactWriteFailure(runDir, result.RunID, err)
	}
	return rowFromRun(runDir, result, contract)
}

func verifierArtifactWriteFailure(runDir, runID string, err error) trace.VerifierResult {
	// verifierArtifactWriteFailure keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return trace.VerifierResult{
		RunID:         runID,
		RunDir:        runDir,
		Result:        trace.VerdictCannotVerify,
		TrustScope:    trace.TrustScopeLocalObserved,
		Completeness:  trace.CompletenessUnknown,
		Replayability: trace.ReplayabilityNone,
		Reason:        fmt.Sprintf("failed writing verifier artifacts: %v", err),
	}
}
