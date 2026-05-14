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
