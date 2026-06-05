package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func writePreviewCommandPayload(stdout io.Writer, payload map[string]any) {
	// Preview output is a declarative plan, not evidence that the command ran.
	data, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintf(stdout, "%s\n", data)
}

func previewCommandPayload(mode string, command []string, contract trace.Contract) map[string]any {
	return map[string]any{
		// The descriptor is the only child-process representation; preview does
		// not execute or retain stdout/stderr.
		"mode":                 mode,
		"command_descriptor":   trace.NewCommandDescriptor(command),
		"contract":             contract,
		"boundaries":           previewBoundaries(),
		"offline_implications": previewOfflineImplications(),
		"writes_artifacts":     false,
		"safe_retention_modes": safeRetentionModes(),
		"warning":              "no run artifacts were written",
	}
}
