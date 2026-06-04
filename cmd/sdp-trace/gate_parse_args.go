package main

import (
	"fmt"
	"io"
)

var gateStringFlags = []string{"out", "contract", "witness", "profile", "checkpoint", "checkpoint-policy"}

func parseGateArgs(args []string, stderr io.Writer) (*flagSet, string, string, int, bool) {
	// Protected-profile inputs are parsed alongside standard gate inputs so one
	// command surface can expose both local and protected modes.
	opts := newStringFlagSet("gate", gateStringFlags)
	if err := opts.parse(args); err != nil {
		// Parse failures happen before any report rows or protected inputs are
		// read.
		fmt.Fprintln(stderr, err)
		return nil, "", "", exitUsage, false
	}
	// Target and output validation stay separate so diagnostics distinguish the
	// evidence source from the artifact destination.
	target, ok := gateTargetArg(opts, stderr)
	if !ok {
		return nil, "", "", exitUsage, false
	}
	outPath, ok := gateOutputPath(opts, stderr)
	if !ok {
		return nil, "", "", exitUsage, false
	}
	return opts, target, outPath, 0, true
}

func gateTargetArg(opts *flagSet, stderr io.Writer) (string, bool) {
	targets := opts.rest()
	if len(targets) == 1 {
		return targets[0], true
	}
	// Gate evaluation is bound to exactly one run root or run directory.
	fmt.Fprintln(stderr, "gate requires <runs-root-or-run-dir>")
	return "", false
}

func gateOutputPath(opts *flagSet, stderr io.Writer) (string, bool) {
	// The output path is validated after target arity so diagnostics first
	// establish which evidence source the gate would evaluate.
	outPath := opts.stringValue("out")
	if outPath != "" {
		return outPath, true
	}
	// Persisted gate JSON is the artifact later explain/preview commands can
	// inspect; stdout is only a rendered copy.
	fmt.Fprintln(stderr, "gate requires --out <file>")
	return "", false
}
