package main

import (
	"context"
	"io"
)

var prReviewHandlers = map[string]subcommandHandler{
	"packet":     runPRReviewPacket,
	"run":        runPRReviewRun,
	"synthesize": runPRReviewSynthesize,
	"validate":   runPRReviewValidate,
	"summarize":  runPRReviewSummarize,
	"check":      runPRReviewCheck,
}

func runPRReview(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runSubcommand(args, stdout, stderr, "pr-review <packet|run|synthesize|validate|summarize|check> [flags]", "pr-review requires packet, run, synthesize, validate, summarize, or check", prReviewHandlers)
}
