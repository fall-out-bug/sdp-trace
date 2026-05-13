package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func resolveProtectedGate(target string, opts *flagSet, stderr io.Writer) (demo.GateResult, int) {
	// Read external trust inputs before loading rows so missing checkpoint,
	// policy, or witness evidence is reported as a gate setup error.
	inputs, code, ok := readProtectedGateInputs(opts, stderr)
	if !ok {
		return demo.GateResult{}, code
	}
	// Rows and contract are loaded after setup evidence so protected gate
	// failures cannot hide missing external authority inputs.
	contract, rows, runDir, code, ok := loadProtectedGateRows(target, opts.stringValue("contract"), stderr)
	if !ok {
		return demo.GateResult{}, code
	}
	// The expected witness is derived from the protected run itself; supplied
	// witness files must match it rather than define their own expectation.
	expected, code, ok := loadProtectedWitnessExpectation(target, stderr)
	if !ok {
		return demo.GateResult{}, code
	}
	checkpointResult := verifiedProtectedCheckpoint(runDir, inputs, expected)
	return evaluateProtectedGate(rows, contract, checkpointResult, inputs.witness, expected), 0
}
