package main

import (
	"context"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

var gateSubcommandHandlers = map[string]subcommandHandler{
	"explain": runGateExplain,
	"preview": runGatePreview,
}

func runGateSubcommand(args []string, stdout, stderr io.Writer) (int, bool) {
	return runOptionalSubcommand(args, stdout, stderr, gateSubcommandHandlers)
}

func runGate(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if code, ok := runGateSubcommand(args, stdout, stderr); ok {
		// Gate subcommands are read-only/explanatory paths and do not evaluate a
		// new gate result.
		return code
	}
	opts, target, outPath, code, ok := parseGateArgs(args, stderr)
	if !ok {
		return code
	}
	if opts.stringValue("profile") == demo.GateProfileProtected {
		// Protected mode requires checkpoint, policy, and witness inputs in
		// addition to local run rows.
		return runProtectedGate(target, outPath, opts, stdout, stderr)
	}
	return runStandardGate(target, outPath, opts, stdout, stderr)
}
