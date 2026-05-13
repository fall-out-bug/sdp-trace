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
