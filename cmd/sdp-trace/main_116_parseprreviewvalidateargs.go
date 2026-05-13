package main

import (
	"fmt"
	"io"
)

func parsePRReviewValidateArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// The validation command requires explicit artifact paths so the resulting
	// JSON can be traced back to concrete packet, profile, run, and ledger files.
	opts := &flagSet{name: "pr-review validate"}
	// Validation output is required so the verdict can be cited by path instead
	// of relying on transient terminal text.
	// Packet/profile/runs/ledger remain independent paths so validation can
	// report malformed or missing evidence at the correct boundary.
	opts.setString("packet", "")
	opts.setString("profile", "")
	opts.setString("runs", "")
	opts.setString("ledger", "")
	// The output path is part of the command contract, not derived from inputs.
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		// Validation cannot begin until every artifact path is parsed.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Rejecting rest arguments before output validation keeps every accepted
	// input represented by a named field in the validation artifact.
	if rejectRest(opts, stderr, "pr-review validate accepts only flags") {
		// Extra positional data would not be represented in validation JSON.
		return nil, exitUsage, false
	}
	// Validate requires a durable output even when the verdict is cannot_verify,
	// because downstream PR gates cite the JSON artifact.
	if err := requireOutputFile("pr-review validate", opts.stringValue("out")); err != nil {
		// Validation output is a machine artifact and must not be implicit.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}
