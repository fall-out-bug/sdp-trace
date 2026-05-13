package main

import (
	"fmt"
	"io"
)

func parsePRReviewSynthesizeArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// Synthesis accepts artifact paths only; there is no inline review payload
	// that could evade packet/run validation.
	opts := &flagSet{name: "pr-review synthesize"}
	// The synthesized ledger is a durable artifact, so the output path is
	// required instead of silently writing only to stdout.
	opts.setString("packet", "")
	opts.setString("runs", "")
	opts.setString("existing-ledger", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		// Bad flags fail before any review artifacts are read.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "pr-review synthesize accepts only flags") {
		// Synthesis inputs are explicit artifact paths only.
		return nil, exitUsage, false
	}
	// The ledger path is validated before artifact reads so a bad output target
	// cannot waste reviewer evidence processing.
	if err := requireOutputFile("pr-review synthesize", opts.stringValue("out")); err != nil {
		// Synthesis output is mandatory because stdout alone is not a stable
		// review ledger reference.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}
