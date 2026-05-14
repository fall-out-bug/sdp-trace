package main

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
	"github.com/fall_out_bug/sdp-trace/internal/telemetry"
)

func renderTelemetryExport(posturePath string) (string, error) {
	var result posture.ExportResult
	if err := readJSONFile(posturePath, &result); err != nil {
		// Missing or malformed posture input means telemetry cannot be verified.
		return "", fmt.Errorf("posture_unreadable")
	}
	rendered, err := telemetry.RenderPrometheus(result)
	if err != nil {
		// Rendering failures preserve cannot_verify instead of emitting partial
		// metrics.
		return "", fmt.Errorf("telemetry_cannot_verify")
	}
	return rendered, nil
}
