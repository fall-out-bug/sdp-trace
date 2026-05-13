package main

import (
	"fmt"
	"io"
)

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
