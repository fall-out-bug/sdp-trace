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
