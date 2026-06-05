package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/query"
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

func writeQueryPackArtifact(opts *queryPackOptions) (int, error) {
	result, err := query.ForensicsBasicPack(opts.runPath)
	if err != nil {
		// Pack generation depends on replayable run evidence.
		return exitCannotVerify, err
	}
	if err := writeJSONFile(opts.outPath, result); err != nil {
		// A generated pack that cannot be persisted is not review evidence.
		return 1, err
	}
	return 0, nil
}
