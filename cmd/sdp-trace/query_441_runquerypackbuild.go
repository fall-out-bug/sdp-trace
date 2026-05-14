package main

import (
	"fmt"
	"io"
)

func runQueryPackBuild(args []string, stderr io.Writer) int {
	opts, err := parseQueryPackArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if err := validateQueryPackOptions(opts); err != nil {
		// Pack/profile validation happens before reading run artifacts so bad
		// command shape cannot be mistaken for unverifiable evidence.
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	// Query-pack build writes a portable JSON artifact for later explanation
	// and review.
	code, err := writeQueryPackArtifact(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return code
	}
	return 0
}
