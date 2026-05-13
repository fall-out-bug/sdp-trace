package main

import "flag"

// optionValues holds raw flag pointers until parsing succeeds.
type optionValues struct {
	// Complexity and MI thresholds are optional gates; zero disables each one.
	cycloOver     *int
	cognitiveOver *int
	miUnder       *float64
	// File and function MI ratchets are kept separate by design.
	fileMIBaseline          *string
	writeFileMIBaseline     *string
	functionMIUnder         *float64
	functionMIBaseline      *string
	writeFunctionMIBaseline *string
	// Output flags change rendering, not measurement.
	gocyclo  *bool
	failOnly *bool
}

// registerFlags declares every supported command-line contract for the tool.
func registerFlags(flags *flag.FlagSet) optionValues {
	// Keep flag registration as the only CLI surface so parseOptions can remain
	// a straight copy from parsed values into immutable run options.
	// Defaults are intentionally non-failing; callers opt into each ratchet or
	// reporting mode with explicit flags.
	return optionValues{
		cycloOver:               flags.Int("cyclo-over", 0, "fail when a production function has cyclomatic complexity greater than this value"),
		cognitiveOver:           flags.Int("cognitive-over", 0, "fail when a production function has cognitive complexity greater than this value"),
		miUnder:                 flags.Float64("mi-under", 0, "fail when a production file has maintainability index below this value"),
		fileMIBaseline:          flags.String("mi-baseline", "", "path to a file maintainability index baseline; files below -mi-under may not regress and new below-threshold files fail"),
		writeFileMIBaseline:     flags.String("write-mi-baseline", "", "write current below-threshold file maintainability index baseline to this path"),
		functionMIUnder:         flags.Float64("function-mi-under", 0, "fail when a production function has maintainability index below this value"),
		functionMIBaseline:      flags.String("function-mi-baseline", "", "path to a function maintainability index baseline; functions below -function-mi-under may not regress and new below-threshold functions fail"),
		writeFunctionMIBaseline: flags.String("write-function-mi-baseline", "", "write current below-threshold function maintainability index baseline to this path"),
		gocyclo:                 flags.Bool("gocyclo", false, "print cyclomatic complexity rows compatible with tools/crapcheck -gocyclo"),
		failOnly:                flags.Bool("fail-only", false, "suppress passing metric rows; threshold and baseline failures still print to stderr"),
	}
}
