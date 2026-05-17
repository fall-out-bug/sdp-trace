package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

// Harness observe is the only harness subcommand that creates a retained
// observation run. Its CLI surface keeps profile, source, and output paths as
// named inputs so downstream validation can cite the artifact boundary instead
// of relying on shell history or prose.

func runHarnessObserve(args []string, stdout, stderr io.Writer) int {
	return runJSONFlagCommand(args, stdout, stderr, parseHarnessObserveArgs, observeHarnessRun, "marshal harness run")
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

func parseHarnessObserveArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// Observation requires the profile, raw source, and output run directory to
	// stay named for replayable artifact provenance.
	return parseFlagOnlyCommand(args, stderr, "harness observe", "harness observe accepts only flags", []observeStringFlag{
		{name: "profile"},
		{name: "source"},
		{name: "out"},
	}, harnessObserveRequiredFlags)
}
