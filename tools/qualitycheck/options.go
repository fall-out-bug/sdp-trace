package main

import (
	"flag"
	"io"
)

// options is the immutable command configuration used after flag parsing.
type options struct {
	// Threshold fields are disabled by zero values.
	cycloOver     int
	cognitiveOver int
	miUnder       float64
	// Baseline paths keep file and function MI ratchets independent.
	fileMIBaseline          string
	writeFileMIBaseline     string
	functionMIUnder         float64
	functionMIBaseline      string
	writeFunctionMIBaseline string
	// Rendering fields affect output, not measured evidence.
	gocyclo  bool
	failOnly bool
	err      io.Writer
}

// parseOptions converts command-line input into options and analysis paths.
func parseOptions(args []string, stderr io.Writer) (options, []string, bool) {
	// Each run owns a FlagSet so tests can exercise parse errors without
	// mutating process-global flags or writing usage to os.Stderr.
	flags := flag.NewFlagSet("qualitycheck", flag.ContinueOnError)
	flags.SetOutput(stderr)
	// Registration stays centralized in option_flags.go; this function only
	// coordinates parse success and path defaulting.
	values := registerFlags(flags)
	if err := flags.Parse(args); err != nil {
		// Stop before path discovery so malformed CLI input cannot accidentally
		// trigger broad repository analysis.
		return options{}, nil, false
	}
	// Positional arguments are analysis scope, while opts carries gate and
	// rendering configuration.
	paths := defaultAnalysisPaths(flags.Args())
	return optionsFromValues(values, stderr), paths, true
}

// optionsFromValues freezes parsed flag pointers into a plain value.
func optionsFromValues(values optionValues, stderr io.Writer) options {
	// Copy parsed pointer values into an immutable plain struct so later
	// analysis, baseline, and report code stay independent from flag.FlagSet.
	return options{
		// Complexity thresholds are AST counts, while MI thresholds are
		// fractional formula results reported to one decimal place.
		cycloOver:     *values.cycloOver,
		cognitiveOver: *values.cognitiveOver,
		miUnder:       *values.miUnder,
		// File and function MI baselines are distinct ratchets; debt in one
		// metric family cannot authorize debt in the other.
		fileMIBaseline:          *values.fileMIBaseline,
		writeFileMIBaseline:     *values.writeFileMIBaseline,
		functionMIUnder:         *values.functionMIUnder,
		functionMIBaseline:      *values.functionMIBaseline,
		writeFunctionMIBaseline: *values.writeFunctionMIBaseline,
		// Output modes affect rendering only; threshold state is still computed
		// from measured evidence later in the run.
		gocyclo:  *values.gocyclo,
		failOnly: *values.failOnly,
		err:      stderr,
	}
}
