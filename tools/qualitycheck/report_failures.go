package main

import "fmt"

// functionFails folds all configured function-level gate failures.
func functionFails(fn functionMetric, opts options, baseline map[string]functionMIBaselineRecord, baselineOK bool) bool {
	// Function gates are independent; evaluate all of them so one failure does
	// not hide a second threshold or MI regression in stderr.
	failed := cyclomaticThresholdFails(fn, opts)
	failed = cognitiveThresholdFails(fn, opts) || failed
	return functionMIThresholdFails(fn, opts, baseline, baselineOK) || failed
}

// cyclomaticThresholdFails evaluates the cyclomatic complexity gate.
func cyclomaticThresholdFails(fn functionMetric, opts options) bool {
	return metricThresholdFails(fn, opts, opts.cycloOver, fn.cyclo, "cyclomatic")
}

// cognitiveThresholdFails evaluates the cognitive complexity gate.
func cognitiveThresholdFails(fn functionMetric, opts options) bool {
	return metricThresholdFails(fn, opts, opts.cognitiveOver, fn.cognitive, "cognitive")
}

// metricThresholdFails applies one integer threshold to one function metric.
func metricThresholdFails(fn functionMetric, opts options, limit int, value int, label string) bool {
	// A non-positive limit disables that gate; this lets one report command mix
	// advisory output with any subset of enforced thresholds.
	if limit <= 0 || value <= limit {
		return false
	}
	fmt.Fprintf(errorOutput(opts), "%s threshold %d exceeded by %s:%d %s\n", label, limit, fn.file, fn.line, fn.name)
	return true
}
