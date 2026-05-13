package main

import (
	"fmt"
	"io"
)

func parseReportArgs(args []string, stderr io.Writer) (*flagSet, string, int, bool) {
	opts := &flagSet{name: "report"}
	// Reports have one durable output root and an optional contract override;
	// the observed run target remains positional for command readability.
	opts.setString("out", "")
	opts.setString("contract", "")
	// Parse flag names before inspecting the target so malformed options are
	// not collapsed into an evidence-root error.
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, "", exitUsage, false
	}
	// A report is useful only when one source and one durable output sink are
	// both explicit in the command line.
	target, ok := reportTargetArg(opts, stderr)
	if !ok {
		return nil, "", exitUsage, false
	}
	if !requireStringFlag(opts, stderr, "out", "report requires --out <dir>") {
		return nil, "", exitUsage, false
	}
	return opts, target, 0, true
}
