package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

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
