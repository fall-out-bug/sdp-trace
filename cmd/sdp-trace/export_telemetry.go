package main

import (
	"fmt"
	"io"
)

func runTelemetryExport(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseTelemetryExportArgs(args, stderr)
	if !ok {
		return code
	}
	// Telemetry is rendered from an already-built posture artifact so the export
	// layer cannot silently broaden repository selection.
	rendered, err := renderTelemetryExport(opts.stringValue("cross-repo-posture"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	if err := writeTelemetryExportOutput(opts.stringValue("out"), rendered, stdout); err != nil {
		// Export errors happen after rendering; no partial metric file is
		// accepted as evidence.
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}
