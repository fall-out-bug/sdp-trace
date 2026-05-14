package main

import (
	"fmt"
	"io"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
)

func runCrossRepoPostureExport(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseCrossRepoPostureExportArgs(args, stderr)
	if !ok {
		return code
	}
	// Build reads the declared selection file and produces one posture artifact;
	// stdout is intentionally ignored so the artifact path stays authoritative.
	result, err := posture.Build(opts.stringValue("selection"), time.Now())
	if err != nil {
		// Build failures are reported without leaking selection parse details
		// into a misleading export artifact.
		fmt.Fprintln(stderr, "no_export_artifact")
		return exitCannotVerify
	}
	_ = stdout
	return writeCrossRepoPostureExport(opts, result, stderr)
}
