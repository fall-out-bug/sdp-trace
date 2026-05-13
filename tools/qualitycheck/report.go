package main

import "io"

// printReport renders measured rows and returns the folded gate verdict.
func printReport(out io.Writer, report qualityReport, opts options) bool {
	// Load function and file MI ratchets independently so a broken optional
	// baseline fails its own gate without suppressing other report output.
	baseline, baselineOK, baselineFailed := loadFunctionBaselineForReport(opts)
	fileBaseline, fileBaselineOK, fileBaselineFailed := loadFileBaselineForReport(opts)
	// Baseline loading is complete before printing begins so output order cannot
	// depend on read failures.
	// Function and file ratchets are joined only at the final boolean gate;
	// report rows remain independently printable.
	// A later failure never suppresses earlier output because this report is the
	// diagnostic surface for ratchet failures.
	// Threshold printers both emit rows and return gate state; keep the final
	// failure fold here so output formatting stays separate from process exit.
	// Use an OR fold instead of early return so users get every failing row in a
	// single run, including baseline read errors and metric threshold breaches.
	failed := baselineFailed
	// Preserve any baseline read failure while continuing to print all metric
	// rows for diagnosis.
	failed = fileBaselineFailed || failed
	// Metric printers are intentionally side-effecting report writers that also
	// return their gate verdict.
	// Function rows are emitted before file rows to match historical baseline
	// review order.
	failed = printFunctionReports(out, report.functions, opts, baseline, baselineOK) || failed
	// File rows use their own baseline type; only the final verdict bool is
	// shared with function-gate output.
	failed = printFileReports(out, report.files, opts, fileBaseline, fileBaselineOK) || failed
	// Returning the folded verdict lets runWithOptions choose the process exit
	// without coupling report formatting to os.Exit.
	return failed
}

// loadFunctionBaselineForReport loads the optional function MI ratchet.
func loadFunctionBaselineForReport(opts options) (map[string]functionMIBaselineRecord, bool, bool) {
	// Function baseline parsing keeps function record shape distinct from file
	// baseline records.
	if missingReportBaseline(opts.functionMIBaseline) {
		return nil, false, false
	}
	baseline, err := readMIBaseline[functionMIBaseline](opts.functionMIBaseline, functionMIBaselineSchema, "function")
	if err != nil {
		return reportBaselineReadError[functionMIBaselineRecord](opts, "function", err)
	}
	return indexFunctionReportBaseline(baseline.Functions), true, false
}

// loadFileBaselineForReport loads the optional file MI ratchet.
func loadFileBaselineForReport(opts options) (map[string]fileMIBaselineRecord, bool, bool) {
	// File baseline parsing is parallel to function parsing but keeps its own
	// schema and record type.
	if missingReportBaseline(opts.fileMIBaseline) {
		return nil, false, false
	}
	baseline, err := readMIBaseline[fileMIBaseline](opts.fileMIBaseline, fileMIBaselineSchema, "file")
	if err != nil {
		return reportBaselineReadError[fileMIBaselineRecord](opts, "file", err)
	}
	return indexFileReportBaseline(baseline.Files), true, false
}
