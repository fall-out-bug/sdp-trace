package main

import (
	"context"
	"fmt"
	"io"
)

func runExport(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if exportTelemetryRequested(args) {
		// Telemetry export consumes posture output and emits Prometheus text.
		return runTelemetryExport(args[1:], stdout, stderr)
	}
	if exportCrossRepoPostureExplainRequested(args) {
		// Cross-repo posture explanations render existing exports only.
		return runCrossRepoPostureExplain(args[2:], stdout, stderr)
	}
	if exportCrossRepoPostureRequested(args) {
		// Cross-repo posture export builds the durable posture artifact.
		return runCrossRepoPostureExport(args[1:], stdout, stderr)
	}
	// Export uses a closed command vocabulary; unsupported exports are usage
	// errors, not unverifiable evidence states.
	fmt.Fprintln(stderr, "export requires cross-repo-posture or telemetry")
	return exitUsage
}
