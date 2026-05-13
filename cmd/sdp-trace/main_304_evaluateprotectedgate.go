package main

import (
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func evaluateProtectedGate(rows []demo.RunRow, contract trace.Contract, checkpointResult checkpoint.VerificationResult, witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) demo.GateResult {
	// Protected evaluation receives already-bound checkpoint and witness facts;
	// the CLI does not override package verdicts.
	return demo.EvaluateProtectedGate(rows, contract, demo.ProtectedGateInput{
		Checkpoint:         checkpointResult,
		PolicyProvided:     true,
		Witness:            &witnessSummary,
		WitnessExpectation: expected,
		Now:                time.Now().UTC(),
	})
}
