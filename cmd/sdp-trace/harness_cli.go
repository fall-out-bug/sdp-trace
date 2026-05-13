package main

import (
	"io"
)

var harnessHandlers = map[string]subcommandHandler{
	"observe":   runHarnessObserve,
	"validate":  runHarnessValidate,
	"summarize": runHarnessSummarize,
}

// The harness command is only a router. Each subcommand owns its evidence
// boundary in a separate file so observe, validate, and summarize cannot blur
// setup artifacts, validation rows, and read-only human summaries.

var harnessObserveRequiredFlags = []requiredCLIFlag{
	{"profile", "harness observe requires --profile"},
	{"source", "harness observe requires --source"},
	{"out", "harness observe requires --out"},
}

var harnessValidateRequiredFlags = []requiredCLIFlag{
	{"profile", "harness validate requires --profile"},
	{"run", "harness validate requires --run"},
}

var harnessSummarizeRequiredFlags = []requiredCLIFlag{
	{"validation", "harness summarize requires --validation"},
}

func runHarness(args []string, stdout, stderr io.Writer) int {
	return runSubcommand(args, stdout, stderr, "harness <observe|validate|summarize> [flags]", "harness requires observe, validate, or summarize", harnessHandlers)
}
