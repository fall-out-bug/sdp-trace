package main

import (
	"fmt"
	"strings"
)

func parseQueryPackArgs(args []string) (*queryPackOptions, error) {
	opts := &flagSet{name: "query-pack"}
	// Pack, run, and output flags are captured before validation so diagnostics
	// can distinguish unsupported pack names from missing evidence paths.
	opts.setString("pack", "")
	opts.setString("run", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		return nil, err
	}
	if len(opts.rest()) != 0 {
		// Pack, run evidence, and output path must be explicit provenance flags.
		// Positional arguments would be omitted from the persisted pack metadata.
		return nil, fmt.Errorf("query-pack accepts only flags")
	}
	return &queryPackOptions{
		// Trimmed values prevent whitespace-only flags from satisfying required
		// evidence-path checks.
		pack:    strings.TrimSpace(opts.stringValue("pack")),
		runPath: strings.TrimSpace(opts.stringValue("run")),
		outPath: strings.TrimSpace(opts.stringValue("out")),
	}, nil
}
