package main

import (
	"fmt"
	"io"
)

var parseObserveSetupCommandArgs = func(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "observe setup"}
	opts.setString("profile", "")
	opts.setString("out", "")
	opts.setString("command", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if !requireOnlyFlags(opts, stderr, "observe setup accepts only flags", observeSetupRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}
