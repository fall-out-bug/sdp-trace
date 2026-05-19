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
	return runHarnessObserveCommand(args, stdout, stderr)
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
	return parseHarnessObserveCommandArgs(args, stderr)
}

func writeHarnessRun(stdout, stderr io.Writer, run harnessobs.Run) int {
	return writeHarnessRunPayload(stdout, stderr, run)
}
