package main

import (
	"context"
	"fmt"
	"io"
)

func runWrap(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	// Legacy wrap keeps parsing separate from recorder execution so malformed
	// wrapper metadata cannot create partial run artifacts.
	opts, command, code, ok := parseWrapArgs(args, stderr)
	if !ok {
		return code
	}
	return runLegacyWrapRecorder(ctx, opts, command, stdout, stderr)
}

func parseWrapArgs(args []string, stderr io.Writer) (*flagSet, []string, int, bool) {
	// Legacy wrap records a command with the default contract unless the caller
	// supplies a contract path.
	opts := &flagSet{name: "wrap"}
	opts.setString("name", "")
	opts.setString("contract", "")
	opts.setString("output-dir", "")
	if err := opts.parse(args); err != nil {
		return nil, nil, exitUsage, false
	}
	// The legacy wrapper still requires an explicit child command; flags only
	// describe recorder metadata.
	command, ok := wrapCommand(opts, stderr)
	if !ok {
		return nil, nil, exitUsage, false
	}
	return opts, command, 0, true
}

func wrapCommand(opts *flagSet, stderr io.Writer) ([]string, bool) {
	command := opts.rest()
	if len(command) == 0 {
		// A wrap without a child command would create an empty provenance shell.
		fmt.Fprintln(stderr, "wrap requires a command")
		return nil, false
	}
	return command, true
}
