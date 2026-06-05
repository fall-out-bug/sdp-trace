package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/interaction"
)

func runInteractionImportTranscript(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseInteractionImportTranscriptArgs(args, stderr)
	if !ok {
		return code
	}
	// Transcript import normalizes external rows into trace events before the CLI
	// emits JSON, keeping summary commands independent of source file shape.
	trace, err := importTranscriptFromOptions(opts)
	return writeImportedTranscript(trace, err, stdout, stderr)
}

func writeImportedTranscript(trace interaction.Trace, err error, stdout, stderr io.Writer) int {
	if err != nil {
		// Import failures mean the transcript source cannot be trusted as trace
		// evidence.
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	writeJSONPayloadUnchecked(stdout, trace)
	return 0
}

func importTranscriptFromOptions(opts *flagSet) (interaction.Trace, error) {
	// Source identity and source ref are preserved so imported transcript events
	// remain attributable after normalization.
	return interaction.ImportTranscript(interaction.ImportOptions{
		TaskID:      opts.stringValue("task-id"),
		Source:      opts.stringValue("source"),
		SourceRef:   opts.stringValue("source-ref"),
		EventsJSONL: opts.stringValue("events-jsonl"),
		Out:         opts.stringValue("out"),
	})
}
