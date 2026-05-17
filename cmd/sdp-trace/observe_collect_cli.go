package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

// Observe collect reads an existing session/run directory and reports evidence.
// It does not execute harness work, infer missing facts from prose, or rewrite
// retained source artifacts.
// The shell status is derived only after the structured payload is emitted.
// Missing or degraded retained evidence remains visible as cannot_verify.

func runObserveCollect(args []string, stdout, stderr io.Writer) int {
	// Collect binds exactly one profile to exactly one retained run directory.
	opts, code, ok := parseObserveCollectArgs(args, stderr)
	if !ok {
		return code
	}
	// The package owns collection state; the CLI preserves and reports it.
	session, observed, err := collectRetainedSession(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return writeObserveCollect(stdout, stderr, session, observed)
}

func collectRetainedSession(opts *flagSet) (harnessobs.SessionRun, harnessobs.Run, error) {
	// Collect replays retained run artifacts and never executes harness work.
	return harnessobs.CollectSession(harnessobs.SessionCollectOptions{
		ProfilePath: opts.stringValue("profile"),
		RunDir:      opts.stringValue("run"),
	})
}

func writeObserveCollect(stdout, stderr io.Writer, session harnessobs.SessionRun, observed harnessobs.Run) int {
	// Emit session and run evidence before translating collection state.
	if !writeObserveRunPayload(stdout, stderr, session, observed, "marshal observe collect") {
		return exitCannotVerify
	}
	return observeCollectExitCode(session)
}

func parseObserveCollectArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// Collection has exactly two provenance inputs: the profile that defines
	// expectations and the run directory that supplies retained evidence.
	return parseFlagOnlyCommand(args, stderr, "observe collect", "observe collect accepts only flags", []observeStringFlag{
		{name: "profile"},
		{name: "run"},
	}, observeCollectRequiredFlags)
}
