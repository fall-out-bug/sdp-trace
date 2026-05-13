package main

import (
	"io"
)

// Observe is the lightweight router for session-observation commands.
// Each handler keeps its own artifact boundary in a sibling file.
// Shared required-flag tables stay here so the top-level command contract is
// visible without mixing setup, collection, and execution logic.

var observeHandlers = map[string]subcommandHandler{
	"setup":   runObserveSetup,
	"collect": runObserveCollect,
	"session": runObserveSession,
}

var observeSetupRequiredFlags = []requiredCLIFlag{
	{"profile", "observe setup requires --profile"},
	{"out", "observe setup requires --out"},
}

var observeCollectRequiredFlags = []requiredCLIFlag{
	{"profile", "observe collect requires --profile"},
	{"run", "observe collect requires --run"},
}

var observeSessionRequiredFlags = []requiredCLIFlag{
	{"profile", "observe session requires --profile"},
	{"out", "observe session requires --out"},
}

func runObserve(args []string, stdout, stderr io.Writer) int {
	return runSubcommand(args, stdout, stderr, "observe <setup|collect|session> [flags]", "observe requires setup, collect, or session", observeHandlers)
}
