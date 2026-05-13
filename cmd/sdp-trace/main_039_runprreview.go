package main

import (
	"context"
	"io"
)

func runPRReview(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runSubcommand(args, stdout, stderr, "pr-review <packet|run|synthesize|validate|summarize|check> [flags]", "pr-review requires packet, run, synthesize, validate, summarize, or check", prReviewHandlers)
}
