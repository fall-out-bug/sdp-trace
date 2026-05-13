package main

import (
	"flag"
	"io"
)

type runOptions struct {
	baseRef   string
	baselines []string
}

func parseRunOptions(args []string, stderr io.Writer) (runOptions, bool) {
	// Build flag state without touching git; option parsing is local evidence.
	flags, baseRef, baselines := newOptionFlags(stderr)
	if err := flags.Parse(args); err != nil {
		// Flag package has already written usage details to stderr.
		return runOptions{}, false
	}
	// Resolution applies policy defaults only after syntax is valid.
	return resolveRunOptions(*baseRef, *baselines, stderr)
}

func newOptionFlags(stderr io.Writer) (*flag.FlagSet, *string, *repeatedFlag) {
	var baselines repeatedFlag
	// This command reads changed paths from stdin; flags only define the base
	// ref and the baseline files that are policy subjects.
	// Parsing is intentionally independent from git so usage errors are local.
	flags := flag.NewFlagSet("mibaselinepolicy", flag.ContinueOnError)
	flags.SetOutput(stderr)
	baseRef := flags.String("base-ref", "", "base git ref used to detect existing MI baselines")
	flags.Var(&baselines, "baseline", "MI baseline path; may be repeated")
	return flags, baseRef, &baselines
}
