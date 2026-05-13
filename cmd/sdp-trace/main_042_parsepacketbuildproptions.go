package main

import (
	"fmt"
	"io"
)

func parsePacketBuildPROptions(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "packet build-pr"}
	// Defaults model the GitHub Actions runtime; explicit flags are used for
	// replay fixtures and local tests.
	opts.setString("source", "github-actions")
	opts.setString("github-event", "")
	opts.setString("checks-json", "")
	opts.setString("artifacts-json", "")
	opts.setString("route-manifest", "")
	opts.setString("github-api-url", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		// Parser errors happen before any review packet or run artifact exists.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Required flag validation keeps the command contract explicit before any
	// GitHub or filesystem evidence is loaded.
	if !requireOnlyFlags(opts, stderr, "packet build-pr accepts only flags", packetBuildPRRequiredFlags) {
		return nil, exitUsage, false
	}
	// The packet builder is flag-only so replay inputs can be reconstructed from
	// invocation text and artifacts.
	return opts, 0, true
}
