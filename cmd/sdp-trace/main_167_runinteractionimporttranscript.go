package main

import (
	"io"
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
