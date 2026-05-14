package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func runPRReviewSynthesize(args []string, stdout, stderr io.Writer) int {
	// Synthesize converts packet/run evidence into a ledger; it never executes
	// reviewers or upgrades review state on its own.
	opts, code, ok := parsePRReviewSynthesizeArgs(args, stderr)
	if !ok {
		return code
	}
	// Synthesis is evidence collation only; unreadable inputs keep the ledger
	// unverifiable rather than producing a partial trust record.
	inputs, err := readPRReviewSynthesisInputs(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	// The ledger is built only after all requested inputs have been decoded, so
	// missing optional history cannot masquerade as a successful merge.
	ledger := prreview.SynthesizeLedger(inputs.packet, inputs.runs, inputs.existing)
	if err := prreview.WriteJSON(opts.stringValue("out"), ledger); err != nil {
		// A synthesized ledger that cannot be written is not durable evidence.
		fmt.Fprintln(stderr, err)
		return 1
	}
	// Stdout mirrors the durable artifact so users inspect the same ledger.
	writeIndentedPayload(stdout, ledger)
	return 0
}
