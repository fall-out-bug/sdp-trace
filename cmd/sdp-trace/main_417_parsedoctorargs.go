package main

import (
	"fmt"
	"io"
)

func parseDoctorArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "doctor"}
	// Doctor defaults inspect local run/report readiness; a profile switches to
	// repository-observer diagnostics.
	opts.setString("contract", "")
	opts.setString("output-dir", defaultRunRoot)
	opts.setString("report-dir", defaultReportDir)
	opts.setString("profile", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		// Doctor parsing is read-only; malformed flags never trigger filesystem
		// diagnostics or repo-observer checks.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if len(opts.rest()) != 0 {
		// Doctor diagnostics are selected only by flags so output remains
		// deterministic and profile-scoped.
		fmt.Fprintln(stderr, "doctor does not accept positional arguments")
		return nil, exitUsage, false
	}
	return opts, 0, true
}
