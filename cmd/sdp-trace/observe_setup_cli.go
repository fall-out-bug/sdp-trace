package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

// Observe setup creates the retained session skeleton only.
// It records profile, output, and an optional command label without executing
// harness work or collecting raw source evidence.
// The JSON response is the durable setup artifact for later collection.
// Any write or marshal failure stays cannot_verify.

func runObserveSetup(args []string, stdout, stderr io.Writer) int {
	// Setup inputs are explicit flags because they define the session artifact.
	opts, code, ok := parseObserveSetupArgs(args, stderr)
	if !ok {
		return code
	}
	// Setup only records metadata; no harness source evidence is collected.
	session, err := setupObservedSession(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return writeObserveSetup(stdout, stderr, session)
}

func setupObservedSession(opts *flagSet) (harnessobs.SessionRun, error) {
	// Setup records metadata only; command is a preview label, not execution.
	return harnessobs.SetupSession(harnessobs.SessionSetupOptions{
		ProfilePath: opts.stringValue("profile"),
		OutDir:      opts.stringValue("out"),
		Command:     opts.stringValue("command"),
	})
}

func writeObserveSetup(stdout, stderr io.Writer, session harnessobs.SessionRun) int {
	// A setup response is evidence only if the CLI can emit structured JSON.
	if !writeJSONPayload(stdout, stderr, session, "marshal observe setup") {
		return exitCannotVerify
	}
	return 0
}

func parseObserveSetupArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "observe setup"}
	// Profile and output identify the session contract and write target;
	// command is only a preview label recorded in setup metadata.
	opts.setString("profile", "")
	opts.setString("out", "")
	opts.setString("command", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Setup is flag-only so profile and output provenance stay explicit.
	if !requireOnlyFlags(opts, stderr, "observe setup accepts only flags", observeSetupRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}
