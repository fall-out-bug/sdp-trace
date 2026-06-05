package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func parsePRReviewCheckArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "pr-review check"}
	registerPRReviewCheckFlags(opts)
	if err := opts.parse(args); err != nil {
		// The combined command still has a flag-only contract.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "pr-review check accepts only flags") {
		// Positional payload would bypass the packet's declared provenance.
		return nil, exitUsage, false
	}
	// Required review anchors are checked after parsing so diagnostics reflect
	// the declared command shape.
	if err := requirePRReviewCheckInputs(opts); err != nil {
		// Missing anchors are caught before any reviewer process can run.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Successful parsing only validates command shape; execution still has to
	// build packet, profile, and run evidence.
	return opts, 0, true
}

func registerPRReviewCheckFlags(opts *flagSet) {
	// Check mode intentionally mirrors packet and run flags so the one-shot path
	// records the same provenance as the decomposed commands.
	registerPRReviewPacketFlags(opts)
	// Profile, runner policy, and work-dir describe the review boundary; preview
	// selects a dry publication path without changing parsed evidence inputs.
	opts.setString("profile", "")
	opts.setString("allow-external-runner", "")
	opts.setString("work-dir", ".")
	opts.setString("not-assessed-reason", "")
	// Preview changes publication only; it does not add evidence inputs.
	opts.setBool("preview", false)
}

func requirePRReviewCheckInputs(opts *flagSet) error {
	outDir := opts.stringValue("out")
	if strings.TrimSpace(outDir) == "" {
		// A combined review check needs a directory because it writes multiple
		// artifacts whose paths become later evidence refs.
		return errors.New("pr-review check requires --out")
	}
	return requirePRReviewPacketInputs(opts)
}
