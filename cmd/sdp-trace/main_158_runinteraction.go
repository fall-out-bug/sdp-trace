package main

import (
	"context"
	"fmt"
	"io"
)

func runInteraction(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	// Interaction commands all materialize or inspect trace artifacts; the
	// router decides only which artifact contract applies.
	if len(args) == 0 {
		// Without a verb there is no interaction evidence contract to apply.
		fmt.Fprintln(stderr, "interaction requires relay, import-transcript, or summarize")
		return exitUsage
	}
	// Interaction commands intentionally share one router so transcript imports,
	// relays, and summaries use the same trace vocabulary.
	handlers := map[string]func([]string, io.Writer, io.Writer) int{
		// Relay needs context because it may execute a forwarded command.
		"relay": func(args []string, stdout, stderr io.Writer) int {
			return runInteractionRelay(ctx, args, stdout, stderr)
		},
		// Import and summarize are pure artifact transforms from the CLI layer.
		"import-transcript": runInteractionImportTranscript,
		"summarize":         runInteractionSummarize,
	}
	handler, ok := handlers[args[0]]
	if !ok {
		// Unknown interaction verbs are command-shape errors, not trace states.
		fmt.Fprintf(stderr, "unknown interaction command: %s\n", args[0])
		return exitUsage
	}
	// The selected handler owns parsing for its artifact shape.
	return handler(args[1:], stdout, stderr)
}
