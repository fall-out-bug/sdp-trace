package main

import (
	"flag"
	"fmt"
	"io"
)

// parseFlags builds the flag set, parses args, and returns the flags.
// A non-negative code means the caller should return immediately.
func parseFlags(args []string, stderr io.Writer) (*flag.FlagSet, *bool, *bool, *string, int) {
	fs := flag.NewFlagSet("osscompat", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = usageFunc(fs, stderr)
	asJSON, list, probe := registerFlags(fs)
	if err := fs.Parse(args); err != nil {
		return fs, asJSON, list, probe, flagParseExitCode(err)
	}
	if len(fs.Args()) > 0 {
		fmt.Fprintf(stderr, "unexpected positional args: %v\n", fs.Args())
		return fs, asJSON, list, probe, 2
	}
	return fs, asJSON, list, probe, -1
}
