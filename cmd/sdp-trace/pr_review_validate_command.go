package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func runPRReviewValidate(args []string, stdout, stderr io.Writer) int {
	// Validate is the CLI gate for review evidence; it lowers unreadable inputs
	// to cannot_verify instead of treating them as a failed review.
	opts, code, ok := parsePRReviewValidateArgs(args, stderr)
	if !ok {
		return code
	}
	// Validation joins independent artifacts and does not trust a ledger unless
	// the packet, profile, and run set can all be loaded.
	inputs, err := readPRReviewValidationInputs(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	// Package validation owns the verdict; the CLI only persists and maps it to
	// process status.
	validation := prreview.Validate(inputs.packet, inputs.profile, inputs.runs, inputs.ledger)
	if err := prreview.WriteJSON(opts.stringValue("out"), validation); err != nil {
		// A validation verdict is useful only after it is persisted.
		fmt.Fprintln(stderr, err)
		return 1
	}
	// The terminal payload is a copy of the persisted validation artifact.
	writeIndentedPayload(stdout, validation)
	return prReviewValidationCLIExitCode(validation)
}

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

func prReviewValidationCLIExitCode(validation prreview.Validation) int {
	if reviewValidationExitCode(validation) != 0 {
		// Invalid review evidence cannot support a PR trust claim.
		return exitCannotVerify
	}
	// A zero exit here only means the review packet validated locally.
	return 0
}
