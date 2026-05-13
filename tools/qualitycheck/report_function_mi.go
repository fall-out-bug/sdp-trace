package main

import "fmt"

// functionMIThresholdFails evaluates the function MI threshold and ratchet.
func functionMIThresholdFails(fn functionMetric, opts options, baseline map[string]functionMIBaselineRecord, baselineOK bool) bool {
	if opts.functionMIUnder <= 0 || fn.maintainabilityIndex >= opts.functionMIUnder {
		return false
	}
	// Without a readable baseline, every below-threshold function is a direct
	// gate failure because there is no ratchet to prove it has not regressed.
	if !baselineOK {
		fmt.Fprintf(errorOutput(opts), "function maintainability index %.1f under threshold %.1f for %s:%d %s\n", fn.maintainabilityIndex, opts.functionMIUnder, fn.file, fn.line, fn.name)
		return true
	}
	return miRecordFails(functionKey(fn), fn.maintainabilityIndex, opts, baseline, "function", "function")
}
