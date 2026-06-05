package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/interaction"
)

func runEnvelope(_ context.Context, args []string, stdout, stderr io.Writer) int {
	// Envelope summarize is a read-only inspection command over one persisted
	// interaction envelope.
	opts, code, ok := parseEnvelopeSummarizeArgs(args, stderr)
	if !ok {
		return code
	}
	// Envelope summaries are inspection artifacts; unreadable envelopes remain
	// cannot_verify rather than being summarized from partial data.
	envelope, err := interaction.ReadEnvelope(opts.stringValue("envelope"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	// Summary generation is pure over the decoded envelope.
	summary := interaction.SummarizeEnvelope(envelope)
	if err := writeOptionalJSONFile(opts.stringValue("out"), summary); err != nil {
		// Envelope summary output is a derived artifact; failed persistence is
		// distinct from envelope readability.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeJSONPayloadUnchecked(stdout, summary)
	return 0
}

func writeOptionalJSONFile(path string, value any) error {
	if strings.TrimSpace(path) == "" {
		// Optional JSON outputs are side effects; an omitted path must not change
		// the command verdict.
		return nil
	}
	return writeJSONFile(path, value)
}
