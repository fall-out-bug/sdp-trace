package main

import (
	"fmt"
	"io"
)

var interactionRelayStringFlags = []struct {
	name         string
	defaultValue string
}{
	{"task-id", ""},
	{"actor-type", "human_user"},
	{"actor-id", ""},
	{"target", "agent"},
	{"event-type", "corrective_feedback"},
	{"operation-id", ""},
	{"stage-id", ""},
	{"out", ""},
}

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

func newInteractionRelayFlagSet() *flagSet {
	opts := &flagSet{name: "interaction relay"}
	// Relay defaults encode a human-to-agent corrective-feedback event; callers
	// override them only when the trace source is more specific.
	for _, flag := range interactionRelayStringFlags {
		opts.setString(flag.name, flag.defaultValue)
	}
	return opts
}

func requireRest(opts *flagSet, stderr io.Writer, message string) bool {
	if len(opts.rest()) != 0 {
		return true
	}
	// Commands after `--` are part of the replay boundary; missing rest args
	// would record feedback without a target command.
	fmt.Fprintln(stderr, message)
	return false
}
