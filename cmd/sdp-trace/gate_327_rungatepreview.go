package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func runGatePreview(args []string, stdout, stderr io.Writer) int {
	opts, targets, code, ok := parseGatePreviewArgs(args, stderr)
	if !ok {
		return code
	}
	if opts.stringValue("profile") == demo.GateProfileProtected {
		// Protected preview reports input readiness only; it never evaluates a
		// protected verdict.
		return runProtectedGatePreview(opts, stdout)
	}
	// Standard preview reads the contract to describe required evidence, not to
	// inspect run artifacts or predict pass/fail.
	contract, err := trace.LoadContract(opts.stringValue("contract"))
	if err != nil {
		// A preview with an unreadable contract cannot describe the expected
		// evidence surface.
		fmt.Fprintln(stderr, err)
		return 1
	}
	// Standard preview emits required evidence only; it never inspects run rows.
	report := buildGatePreviewReport(contract, opts.stringValue("witness"), targets[0])
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}
