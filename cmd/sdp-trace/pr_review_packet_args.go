package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
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

func registerPRReviewPacketFlags(opts *flagSet) {
	// Packet metadata is fully flag-driven so generated review packets can be
	// replayed without hidden process context.
	for _, flag := range prReviewPacketStringFlags {
		opts.setString(flag.name, flag.defaultValue)
	}
}

func requirePRReviewPacketInputs(opts *flagSet) error {
	for _, flag := range prReviewPacketRequiredFlags {
		if strings.TrimSpace(opts.stringValue(flag.name)) == "" {
			// Required packet fields are provenance anchors, not cosmetic labels.
			return errors.New(flag.message)
		}
	}
	return nil
}
