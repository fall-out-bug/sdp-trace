package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

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
	opts := &flagSet{name: "observe collect"}
	// Collection has exactly two provenance inputs: the profile that defines
	// expectations and the run directory that supplies retained evidence.
	opts.setString("profile", "")
	opts.setString("run", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Collection inputs are named flags because the run directory is evidence.
	if !requireOnlyFlags(opts, stderr, "observe collect accepts only flags", observeCollectRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func writeObserveRunPayload(stdout, stderr io.Writer, session harnessobs.SessionRun, observed harnessobs.Run, message string) bool {
	// Emit session metadata beside observed run evidence for replay/debugging.
	payload := struct {
		Session harnessobs.SessionRun `json:"session"`
		Run     harnessobs.Run        `json:"run"`
	}{Session: session, Run: observed}
	return writeJSONPayload(stdout, stderr, payload, message)
}

func observeCollectExitCode(session harnessobs.SessionRun) int {
	if session.CollectionState == harnessobs.StateCannotVerify {
		// Collection can write a session payload while still failing as evidence.
		return exitCannotVerify
	}
	return 0
}

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
