package main

import (
	"fmt"
	"io"
)

func parseInteractionRelayArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// Relay parsing keeps trace identity flags separate from the command after
	// `--`, which is forwarded and recorded exactly.
	opts := newInteractionRelayFlagSet()
	if err := opts.parse(args); err != nil {
		// Flag parse errors happen before the command boundary can be trusted.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// The task and output path are the minimum durable trace coordinates.
	if !requireRequiredFlags(opts, stderr, interactionRelayRequiredFlags) {
		return nil, exitUsage, false
	}
	// Rest arguments are mandatory because relay exists to bind feedback to a
	// concrete forwarded command.
	if !requireRest(opts, stderr, "interaction relay requires forward command after --") {
		return nil, exitUsage, false
	}
	// The remaining args after `--` are forwarded exactly by Relay and become
	// part of the recorded command boundary.
	return opts, 0, true
}
