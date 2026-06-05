package main

import (
	"fmt"
	"io"
)

func parseReleaseProofArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "release-proof"}
	// The default manifest is an example contract; callers still need an output
	// path so the generated proof is a durable artifact.
	opts.setString("manifest", "examples/contract-foundation/contract-manifest.example.json")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if !requireOnlyFlags(opts, stderr, "release-proof accepts only flags", releaseProofRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}
