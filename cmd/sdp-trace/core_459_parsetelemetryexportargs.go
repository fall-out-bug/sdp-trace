package main

import (
	"fmt"
	"io"
)

func parseTelemetryExportArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "export telemetry"}
	// Telemetry export accepts only the renderer profile, posture artifact, and
	// output target needed to replay metric generation.
	opts.setString("profile", "")
	opts.setString("cross-repo-posture", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Profile, posture artifact, and destination stay named so exports are
	// auditable from the command line alone.
	if rejectRest(opts, stderr, "export telemetry accepts only flags") {
		return nil, exitUsage, false
	}
	return requireTelemetryExportArgs(opts, stderr)
}
