package main

import "fmt"

// miBaselineRecord is the minimal comparison contract shared by MI ratchets.
type miBaselineRecord interface {
	baselineMI() float64
}

// miRecordFails evaluates whether one below-threshold MI row is authorized by
// its baseline.
func miRecordFails[R miBaselineRecord](key string, currentMI float64, opts options, baseline map[string]R, label string, subject string) bool {
	record, ok := baseline[key]
	if !ok {
		fmt.Fprintf(errorOutput(opts), "%s MI baseline missing for below-threshold %s %s\n", label, subject, key)
		return true
	}
	// Baseline comparison uses the same one-decimal value persisted to disk so
	// report verdicts do not depend on hidden floating-point precision.
	current := roundMetric(currentMI)
	if current < record.baselineMI() {
		fmt.Fprintf(errorOutput(opts), "%s MI baseline regressed for %s: %.1f < %.1f\n", label, key, current, record.baselineMI())
		return true
	}
	return false
}

// baselineMI exposes the persisted function MI value for generic comparison.
func (record functionMIBaselineRecord) baselineMI() float64 {
	return record.MaintainabilityIndex
}

// baselineMI exposes the persisted file MI value for generic comparison.
func (record fileMIBaselineRecord) baselineMI() float64 {
	return record.MaintainabilityIndex
}
