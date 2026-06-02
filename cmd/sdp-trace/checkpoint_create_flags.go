package main

import (
	"fmt"
	"io"
)

func parseCheckpointCreateArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := newCheckpointCreateFlagSet()
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Required create flags provide all source, sink, and signer inputs for a
	// replayable checkpoint artifact.
	if !requireOnlyFlags(opts, stderr, "checkpoint create accepts only flags", checkpointCreateRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func newCheckpointCreateFlagSet() *flagSet {
	opts := &flagSet{name: "checkpoint create"}
	for _, flag := range checkpointCreateStringFlags {
		// Registration order is fixed so help/tests observe the same command
		// contract while defaults stay beside their flag names.
		opts.setString(flag.name, flag.defaultValue)
	}
	return opts
}
