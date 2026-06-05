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

func dispatchSubcommand(cmd string, args []string, stdout, stderr io.Writer, label string, handlers map[string]subcommandHandler) int {
	handler, ok := handlers[cmd]
	if !ok {
		// Keep command-family names stable in diagnostics even when the usage label
		// contains argument suffixes.
		fmt.Fprintf(stderr, "unknown %s command: %s\n", subcommandName(label), cmd)
		return exitUsage
	}
	return handler(args, stdout, stderr)
}
