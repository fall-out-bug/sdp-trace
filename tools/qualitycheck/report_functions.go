package main

import (
	"fmt"
	"io"
)

// printFunctionReports writes function rows and folds their gate verdicts.
func printFunctionReports(out io.Writer, functions []functionMetric, opts options, baseline map[string]functionMIBaselineRecord, baselineOK bool) bool {
	failed := false
	for _, fn := range functions {
		// Emit each row before checking gates so passing output order and failure
		// diagnostics stay stable across threshold combinations.
		printFunctionMetric(out, fn, opts)
		failed = functionFails(fn, opts, baseline, baselineOK) || failed
	}
	return failed
}

// printFunctionMetric writes one function row for the selected output mode.
func printFunctionMetric(out io.Writer, fn functionMetric, opts options) {
	if opts.gocyclo {
		// Gocyclo compatibility is an alternate stdout contract, not an extra
		// column set on the default qualitycheck report.
		fmt.Fprintf(out, "%d %s %s %s:%d:%d\n", fn.cyclo, fn.packageName, fn.name, fn.file, fn.line, fn.column)
		return
	}
	if opts.failOnly {
		return
	}
	fmt.Fprintf(out, "%s:%d %s cyclo=%d cognitive=%d maintainability_index=%.1f\n", fn.file, fn.line, fn.name, fn.cyclo, fn.cognitive, fn.maintainabilityIndex)
}
