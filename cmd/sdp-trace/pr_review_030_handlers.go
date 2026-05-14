package main

var prReviewHandlers = map[string]subcommandHandler{
	"packet":     runPRReviewPacket,
	"run":        runPRReviewRun,
	"synthesize": runPRReviewSynthesize,
	"validate":   runPRReviewValidate,
	"summarize":  runPRReviewSummarize,
	"check":      runPRReviewCheck,
}
