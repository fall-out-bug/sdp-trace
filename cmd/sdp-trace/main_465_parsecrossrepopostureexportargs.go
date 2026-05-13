package main

import (
	"fmt"
	"io"
)

func parseCrossRepoPostureExportArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "export cross-repo-posture"}
	// Profile names the posture contract, selection names the repository set,
	// and out names the durable result artifact.
	opts.setString("profile", "")
	opts.setString("selection", "")
	opts.setString("out", "")
	opts.setBool("validate-only", false)
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Selection, profile, and output remain explicit flags for replayable export
	// provenance.
	if rejectRest(opts, stderr, "export cross-repo-posture accepts only flags") {
		return nil, exitUsage, false
	}
	return requireCrossRepoPostureExportArgs(opts, stderr)
}
