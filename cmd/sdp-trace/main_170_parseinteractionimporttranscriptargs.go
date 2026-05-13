package main

import (
	"fmt"
	"io"
)

func parseInteractionImportTranscriptArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// Import parsing requires the external transcript path and the target trace
	// output path up front.
	opts := &flagSet{name: "interaction import-transcript"}
	// Imported transcripts require both input and output paths to avoid
	// terminal-only trace evidence.
	opts.setString("source", "")
	opts.setString("source-ref", "")
	opts.setString("task-id", "")
	opts.setString("events-jsonl", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		// Malformed flags stop before transcript rows are read.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if !requireOnlyFlags(opts, stderr, "interaction import-transcript accepts only flags", interactionImportTranscriptRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}
