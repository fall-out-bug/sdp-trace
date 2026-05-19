package main

import (
	"fmt"
	"io"
)

var parseHarnessObserveCommandArgs = func(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "harness observe"}
	opts.setString("profile", "")
	opts.setString("source", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if !requireOnlyFlags(opts, stderr, "harness observe accepts only flags", harnessObserveRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}
