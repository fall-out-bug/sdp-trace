package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

// Observe session is the combined setup, execution, and collection path.
// The command after -- is the only child process payload; observer flags remain
// separate so child flags cannot be mistaken for provenance inputs.
// The emitted payload matches the collect shape for reviewer replay.

func runObserveSession(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseObserveSessionArgs(args, stderr)
	if !ok {
		return code
	}
	// Session mode executes and collects in one call while preserving the same
	// combined output shape as setup+collect.
	session, observed, code, ok := collectObservedSession(opts, stderr)
	if !ok {
		return code
	}
	if !writeObserveRunPayload(stdout, stderr, session, observed, "marshal observe session") {
		return exitCannotVerify
	}
	return 0
}

func collectObservedSession(opts *flagSet, stderr io.Writer) (harnessobs.SessionRun, harnessobs.Run, int, bool) {
	// The command payload after -- is forwarded to harnessobs as observed work.
	// Parse-time flags remain separate from the command argv so child flags are
	// not mistaken for observer options.
	session, observed, err := harnessobs.RunSession(harnessobs.SessionOptions{
		ProfilePath: opts.stringValue("profile"),
		OutDir:      opts.stringValue("out"),
		Command:     opts.rest(),
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return harnessobs.SessionRun{}, harnessobs.Run{}, exitCannotVerify, false
	}
	return session, observed, 0, true
}

func parseObserveSessionArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "observe session"}
	// Session mode combines setup and collection but keeps the same explicit
	// profile/output inputs as setup.
	opts.setString("profile", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if !requireRequiredFlags(opts, stderr, observeSessionRequiredFlags) {
		return nil, exitUsage, false
	}
	// After required flags, the remaining argv must be the observed command.
	if len(opts.rest()) == 0 {
		// Without a child command, session mode would create setup metadata with
		// no observed execution evidence.
		// Session observation needs a command payload to produce replayable run
		// evidence.
		fmt.Fprintln(stderr, "observe session requires command after --")
		return nil, exitUsage, false
	}
	return opts, 0, true
}
