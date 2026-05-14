package main

import (
	"fmt"
	"io"
)

func runSubcommand(args []string, stdout, stderr io.Writer, label, usage string, handlers map[string]subcommandHandler) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, usage)
		return exitUsage
	}
	if isHelpArg(args[0]) {
		// Help for a command family is local CLI documentation, not an evidence
		// producing operation.
		fmt.Fprintf(stdout, "Usage: sdp-trace %s\n", label)
		return 0
	}
	return dispatchSubcommand(args[0], args[1:], stdout, stderr, label, handlers)
}
