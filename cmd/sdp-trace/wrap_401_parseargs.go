package main

import (
	"io"
)

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
