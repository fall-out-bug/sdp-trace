package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func loadProtectedGateRows(target, contractPath string, stderr io.Writer) (trace.Contract, []demo.RunRow, string, int, bool) {
	// Protected rows are always checked against the requested contract; this
	// path has no silent default upgrade after the profile is selected.
	contract, err := trace.LoadContract(contractPath)
	if err != nil {
		// Protected mode cannot fall back to an implicit contract after the user
		// supplied one.
		fmt.Fprintln(stderr, err)
		return trace.Contract{}, nil, "", 1, false
	}
	rows, err := demo.VerifiedRows(target, contract)
	if err != nil {
		// Row replay failures block protected evaluation before checkpoint facts
		// can be joined.
		fmt.Fprintln(stderr, err)
		return trace.Contract{}, nil, "", 1, false
	}
	// Protected checkpoints replay a single concrete run directory, not only a
	// report-level collection of rows.
	runDir, err := protectedRunDir(target)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return trace.Contract{}, nil, "", exitCannotVerify, false
	}
	return contract, rows, runDir, 0, true
}

func loadProtectedWitnessExpectation(target string, stderr io.Writer) (demo.WitnessExpectation, int, bool) {
	// The expectation loader derives the run id and artifact digests that the
	// supplied witness summary must bind to.
	expected, err := demoWitnessExpectation(target)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return demo.WitnessExpectation{}, exitCannotVerify, false
	}
	return expected, 0, true
}
