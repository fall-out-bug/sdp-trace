package main

import (
	"fmt"
	"io"
)

func parsePRReviewPacketArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "pr-review packet"}
	registerPRReviewPacketFlags(opts)
	if err := opts.parse(args); err != nil {
		// Parser errors are command-shape failures before any packet evidence is
		// copied or hashed.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "pr-review packet accepts only flags") {
		// Positional arguments would be hidden packet inputs, so reject them.
		return nil, exitUsage, false
	}
	if err := requirePRReviewPacketInputs(opts); err != nil {
		// Missing packet anchors are usage errors because the packet cannot be
		// constructed at all.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}
