package main

import (
	"fmt"
	"io"
	"os"
)

// main exits with the qualitycheck command verdict.
func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run parses CLI options and dispatches the normalized request.
func run(args []string, stdout, stderr io.Writer) int {
	// Argument parsing is the CLI boundary; the rest of the package receives a
	// normalized options struct and explicit writers.
	opts, paths, ok := parseOptions(args, stderr)
	if !ok {
		return 2
	}
	return runWithOptions(paths, stdout, stderr, opts)
}

// runWithOptions measures paths, optionally writes baselines, and returns the
// gate exit code.
func runWithOptions(paths []string, stdout, stderr io.Writer, opts options) int {
	report, err := analyzePaths(paths)
	if err != nil {
		fmt.Fprintf(stderr, "analyze paths: %v\n", err)
		return 2
	}
	// Baseline writes are a terminal mode: callers either refresh ratchets or
	// evaluate gates, never both in the same invocation.
	if code, ok := writeRequestedBaseline(report, stderr, opts); ok {
		return code
	}
	// Reporting owns the pass/fail boundary after discovery and measurement
	// have produced a complete in-memory quality report.
	if printReport(stdout, report, opts) {
		return 1
	}
	return 0
}

// writeRequestedBaseline handles terminal baseline-refresh modes.
func writeRequestedBaseline(report qualityReport, stderr io.Writer, opts options) (int, bool) {
	// Baseline flags are mutually terminal modes; the bool tells runWithOptions
	// that normal threshold reporting has already been bypassed.
	if opts.writeFunctionMIBaseline != "" {
		return writeBaselineExit(report, stderr, opts.writeFunctionMIBaseline, opts.functionMIUnder, "function", writeFunctionMIBaselineFile), true
	}
	if opts.writeFileMIBaseline != "" {
		return writeBaselineExit(report, stderr, opts.writeFileMIBaseline, opts.miUnder, "file", writeFileMIBaselineFile), true
	}
	return 0, false
}
