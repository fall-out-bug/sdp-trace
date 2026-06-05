package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func runProtectedGate(target, outPath string, opts *flagSet, stdout, stderr io.Writer) int {
	// Protected gate resolution is separated from writing so input failures do
	// not create a partial gate artifact.
	result, code := resolveProtectedGate(target, opts, stderr)
	if code != 0 {
		return code
	}
	return writeProtectedGateResult(outPath, result, stdout, stderr)
}

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

func writeProtectedGateResult(path string, result demo.GateResult, stdout, stderr io.Writer) int {
	if err := writeJSONFile(path, result); err != nil {
		// A protected verdict that cannot be written is not reviewable evidence.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeIndentedPayload(stdout, result)
	return gateExitCode(result)
}
