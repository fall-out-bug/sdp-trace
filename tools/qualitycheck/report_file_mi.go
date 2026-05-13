package main

import "fmt"

// fileFails applies the file MI threshold, including optional baseline ratchets.
func fileFails(file fileMetric, opts options, baseline map[string]fileMIBaselineRecord, baselineOK bool) bool {
	if opts.miUnder <= 0 || file.maintainabilityIndex >= opts.miUnder {
		return false
	}
	// A readable file baseline turns below-threshold files into ratcheted debt;
	// without it, the raw threshold breach is the gate failure.
	if baselineOK {
		return miRecordFails(fileKey(file), file.maintainabilityIndex, opts, baseline, "file", "file")
	}
	fmt.Fprintf(errorOutput(opts), "maintainability index %.1f under threshold %.1f for %s\n", file.maintainabilityIndex, opts.miUnder, file.file)
	return true
}
