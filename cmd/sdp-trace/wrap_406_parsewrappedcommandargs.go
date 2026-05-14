package main

import (
	"io"
)

func parseWrappedCommandArgs(args []string, stderr io.Writer) (*flagSet, []string, int, bool) {
	opts := &flagSet{name: "run"}
	// The modern run command keeps task identity and contract selection explicit
	// before the observed command payload.
	opts.setString("task", "")
	opts.setString("contract", "")
	opts.setBool("use-default-contract", false)
	opts.setString("name", "")
	opts.setString("output-dir", "")
	if err := opts.parse(args); err != nil {
		// Run flag parsing fails before command execution, so no partial trace is
		// created for malformed recorder options.
		return nil, nil, exitUsage, false
	}
	command := opts.rest()
	// The remaining argv is the command under observation; flags after this
	// point belong to the child process.
	if !requireWrappedCommandArgs(opts, command, stderr) {
		return nil, nil, exitUsage, false
	}
	return opts, command, 0, true
}
