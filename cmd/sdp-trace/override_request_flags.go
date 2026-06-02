package main

import (
	"fmt"
	"io"
)

func parseOverrideRequestArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	if !isOverrideRequest(args) {
		// The override namespace currently has one explicit write action.
		fmt.Fprintln(stderr, "override requires request")
		return nil, exitUsage, false
	}
	return parseOverrideRequestFlags(args[1:], stderr)
}

func isOverrideRequest(args []string) bool {
	return len(args) != 0 && args[0] == "request"
}

func parseOverrideRequestFlags(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// Override requests are write operations, so parsing must establish a fully
	// named payload before any trace event can be appended.
	opts := &flagSet{name: "override request"}
	// Each flag maps directly to a persisted trace payload key, keeping the
	// event reviewable without positional inference.
	// Requiredness is checked after parsing so diagnostics can distinguish
	// unknown flags from missing trace fields.
	opts.setString("out", "")
	opts.setString("id", "")
	opts.setString("by", "")
	opts.setString("reason", "")
	opts.setString("source-ref", "")
	opts.setString("scope", "")
	opts.setString("external-reference", "")
	if err := opts.parse(args); err != nil {
		// Parser errors stop before the command can create partial override
		// evidence.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// No positional text is accepted because all persisted override request
	// fields must have stable JSON keys.
	if len(opts.rest()) != 0 {
		// Free-form positional text would make the override reason ambiguous.
		fmt.Fprintln(stderr, "override request accepts only flags")
		return nil, exitUsage, false
	}
	return requireOverrideRequestFlags(opts, stderr)
}

func requireOverrideRequestFlags(opts *flagSet, stderr io.Writer) (*flagSet, int, bool) {
	// Required field validation happens before the run directory is opened for
	// append, preventing partial override events.
	if !requireRequiredFlags(opts, stderr, overrideRequestRequiredFlags) {
		// Required fields identify who asked, what scope is affected, and which
		// source evidence the override references.
		return nil, exitUsage, false
	}
	return opts, 0, true
}

var overrideRequestRequiredFlags = []requiredCLIFlag{
	{"out", "override request requires --out"},
	{"id", "override request requires --id"},
	{"by", "override request requires --by"},
	{"reason", "override request requires --reason"},
	{"source-ref", "override request requires --source-ref"},
	{"scope", "override request requires --scope"},
}
