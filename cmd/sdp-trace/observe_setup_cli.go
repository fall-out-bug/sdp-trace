package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

// Observe setup creates the retained session skeleton only.
// It records profile, output, and an optional command label without executing
// harness work or collecting raw source evidence.
// The JSON response is the durable setup artifact for later collection.
// Any write or marshal failure stays cannot_verify.

func runObserveSetup(args []string, stdout, stderr io.Writer) int {
	return runObserveSetupCommand(args, stdout, stderr)
}

func setupObservedSession(opts *flagSet) (harnessobs.SessionRun, error) {
	// Setup records metadata only; command is a preview label, not execution.
	return harnessobs.SetupSession(harnessobs.SessionSetupOptions{
		ProfilePath: opts.stringValue("profile"),
		OutDir:      opts.stringValue("out"),
		Command:     opts.stringValue("command"),
	})
}

func parseObserveSetupArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	return parseObserveSetupCommandArgs(args, stderr)
}

func writeObserveSetup(stdout, stderr io.Writer, session harnessobs.SessionRun) int {
	return writeObserveSetupPayload(stdout, stderr, session)
}
