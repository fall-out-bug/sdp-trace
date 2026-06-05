package main

import (
	"fmt"
	"io"
	"strings"
)

var checkpointVerifyStringFlags = []string{"run", "checkpoint", "policy"}

func parseCheckpointVerifyArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := newCheckpointVerifyFlagSet()
	if err := opts.parse(args); err != nil {
		// Parse errors happen before any signed checkpoint or policy is loaded.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if len(opts.rest()) != 0 {
		// Positional arguments would make the verification source ambiguous.
		fmt.Fprintln(stderr, "checkpoint verify accepts only flags")
		return nil, exitUsage, false
	}
	if err := requireCheckpointVerifyInputs(opts); err != nil {
		// Required verification inputs identify the run replay source and signed
		// checkpoint artifact before optional policy can affect trust scope.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func newCheckpointVerifyFlagSet() *flagSet {
	opts := &flagSet{name: "checkpoint verify"}
	// Verification names the run, signed checkpoint, and optional trust policy
	// as separate inputs so each can fail independently.
	for _, flag := range checkpointVerifyStringFlags {
		opts.setString(flag, "")
	}
	return opts
}

func requireCheckpointVerifyInputs(opts *flagSet) error {
	if strings.TrimSpace(opts.stringValue("run")) == "" {
		// The run directory is the source replay target for the signed payload.
		return fmt.Errorf("checkpoint verify requires --run")
	}
	if strings.TrimSpace(opts.stringValue("checkpoint")) == "" {
		// The signed checkpoint artifact is mandatory verification input.
		return fmt.Errorf("checkpoint verify requires --checkpoint")
	}
	return nil
}
