package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

// Harness observe is the only harness subcommand that creates a retained
// observation run. Its CLI surface keeps profile, source, and output paths as
// named inputs so downstream validation can cite the artifact boundary instead
// of relying on shell history or prose.

func runHarnessObserve(args []string, stdout, stderr io.Writer) int {
	// Parse-time validation keeps every evidence input named as a flag.
	opts, code, ok := parseHarnessObserveArgs(args, stderr)
	if !ok {
		return code
	}
	// Observation work stays in harnessobs; the CLI only reports its artifact.
	run, err := observeHarnessRun(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return writeHarnessRun(stdout, stderr, run)
}

func observeHarnessRun(opts *flagSet) (harnessobs.Run, error) {
	// The CLI preserves the observation boundary: source parsing,
	// normalization, and verdict derivation stay package-owned.
	return harnessobs.Observe(harnessobs.ObserveOptions{
		ProfilePath: opts.stringValue("profile"),
		SourcePath:  opts.stringValue("source"),
		OutDir:      opts.stringValue("out"),
	})
}

func writeHarnessRun(stdout, stderr io.Writer, run harnessobs.Run) int {
	// A marshal failure means the command cannot publish replayable evidence.
	if !writeJSONPayload(stdout, stderr, run, "marshal harness run") {
		return exitCannotVerify
	}
	return 0
}

func parseHarnessObserveArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "harness observe"}
	// Observation requires the profile, raw source, and output run directory to
	// stay named for replayable artifact provenance.
	opts.setString("profile", "")
	opts.setString("source", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Harness observation is flag-only so profile/source/output are auditable.
	if !requireOnlyFlags(opts, stderr, "harness observe accepts only flags", harnessObserveRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}
