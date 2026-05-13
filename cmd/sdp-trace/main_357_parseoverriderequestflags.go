package main

import (
	"fmt"
	"io"
)

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
