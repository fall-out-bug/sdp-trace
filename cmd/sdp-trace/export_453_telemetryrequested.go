package main

func exportTelemetryRequested(args []string) bool {
	return exportCommandIs(args, "telemetry")
}
