package main

import (
	"fmt"
	"io"
)

// baselineWriter is the common writer shape for function and file MI ratchets.
type baselineWriter func(string, qualityReport, float64) error

// writeBaselineExit converts a baseline write result into a CLI exit code.
func writeBaselineExit(report qualityReport, stderr io.Writer, path string, threshold float64, label string, write baselineWriter) int {
	// Keep writer selection injected so function and file baselines share one
	// exit-code path without hiding the target-specific serializer.
	if err := write(path, report, threshold); err != nil {
		fmt.Fprintf(stderr, "write %s MI baseline: %v\n", label, err)
		return 2
	}
	return 0
}
