package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
	"path/filepath"
)

func rowFromRun(runDir string, result trace.VerifierResult, contract trace.Contract) RunRow {
	// rowFromRun keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	row := rowFromVerifierResult(runDir, result)
	artifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {

		row.Kind = "unmatched"
		row.KindReason = "run artifact could not be loaded"
		return row
	}
	applyRunArtifact(&row, artifact, contract)
	return row
}
func rowFromVerifierResult(runDir string, result trace.VerifierResult) RunRow {
	// rowFromVerifierResult keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return RunRow{
		Name:          filepath.Base(runDir),
		RunID:         result.RunID,
		Result:        result.Result,
		TrustScope:    result.TrustScope,
		Completeness:  result.Completeness,
		Replayability: result.Replayability,
		Reason:        result.Reason,

		ClosureState: trace.ClosureStateUnknown,
	}
}
