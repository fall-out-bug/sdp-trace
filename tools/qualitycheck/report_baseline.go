package main

import (
	"fmt"
	"strings"
)

// missingReportBaseline reports whether a baseline flag requested no file.
func missingReportBaseline(path string) bool {
	// Empty path means no ratchet was requested; threshold-only reporting can
	// continue without marking baseline state as failed.
	return strings.TrimSpace(path) == ""
}

// reportBaselineReadError converts a baseline read failure into gate state.
func reportBaselineReadError(opts options, label string, err error) (bool, bool) {
	// Keep the baseline error attached to the metric family label so users can
	// distinguish function and file ratchet failures.
	fmt.Fprintf(errorOutput(opts), "read %s MI baseline: %v\n", label, err)
	return false, true
}

// indexFunctionReportBaseline builds function-ratchet lookups by stable key.
func indexFunctionReportBaseline(items []functionMIBaselineRecord) map[string]functionMIBaselineRecord {
	return indexReportBaseline(items, func(record functionMIBaselineRecord) string { return record.Key })
}

// indexFileReportBaseline builds file-ratchet lookups by stable key.
func indexFileReportBaseline(items []fileMIBaselineRecord) map[string]fileMIBaselineRecord {
	return indexReportBaseline(items, func(record fileMIBaselineRecord) string { return record.Key })
}

// indexReportBaseline converts persisted baseline rows into a lookup map.
func indexReportBaseline[R any](items []R, key func(R) string) map[string]R {
	// Persisted record order is irrelevant after loading; report order is owned
	// by current metric measurement.
	// The map is a lookup cache only; fresh metric order still controls report
	// order so output stays deterministic.
	byKey := make(map[string]R, len(items))
	for _, record := range items {
		// Duplicate detection belongs to baseline validation. Loading preserves
		// normal map semantics by letting the later row win.
		// The reporter only needs the final lookup value for each key.
		byKey[key(record)] = record
	}
	return byKey
}
