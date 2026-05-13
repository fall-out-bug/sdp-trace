package main

import (
	"fmt"
	"io"
)

// printFileReports writes file rows and folds their gate verdicts.
func printFileReports(out io.Writer, files []fileMetric, opts options, baseline map[string]fileMIBaselineRecord, baselineOK bool) bool {
	failed := false
	if opts.gocyclo {
		// Gocyclo mode is function-only output; file MI still stays available
		// through the normal report mode.
		return false
	}
	for _, file := range files {
		printFileMetric(out, file, opts)
		if fileFails(file, opts, baseline, baselineOK) {
			failed = true
		}
	}
	return failed
}

// printFileMetric writes one whole-file MI row when output is not fail-only.
func printFileMetric(out io.Writer, file fileMetric, opts options) {
	// fail-only suppresses advisory rows but never suppresses the gate checks
	// that run in printFileReports.
	if opts.failOnly {
		return
	}
	fmt.Fprintf(out, "%s maintainability_index=%.1f lines=%d cyclo=%d halstead_volume=%.1f\n", file.file, file.maintainabilityIndex, file.lines, file.cyclo, file.halsteadVolume)
}
