package main

import (
	"fmt"
	"io"
)

// The report has two trust surfaces:
// stdout is the complete evidence table for replay and review,
// stderr is only the failing-gate summary,
// and the integer return is the machine verdict consumed by callers.
// Keep those surfaces independent so log truncation cannot hide evidence rows.
// A row that passes the threshold is still printed because absence is not proof.
// The failed slice records only threshold breaches; it must not filter stdout.
// That makes the emitted table stable across passing and failing gate runs.
// Stderr deliberately names the threshold once rather than replaying rows.
// Callers can archive stdout as evidence even when the gate returns zero.
func printResults(stdout io.Writer, stderr io.Writer, rows []resultRow, threshold float64, strictLess bool) int {
	var failed []resultRow
	for _, row := range rows {
		// Always print the evidence row; callers decide whether stdout is a
		// human report or a captured artifact for later review.
		fmt.Fprintf(stdout, "%s:%s %s complexity=%d coverage=%.1f crap=%.2f\n", row.file, row.line, row.function, row.complexity, row.coverage, row.crap)
		if exceedsThreshold(row.crap, threshold, strictLess) {
			failed = append(failed, row)
		}
	}
	// The exit code is a gate verdict, while stderr carries the compact failure
	// summary used by CI logs.
	if len(failed) > 0 {
		fmt.Fprintf(stderr, "CRAP threshold %.2f exceeded by %d function(s)\n", threshold, len(failed))
		return 1
	}
	return 0
}
