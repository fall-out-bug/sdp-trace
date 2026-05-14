package main

import (
	"fmt"
	"io"
)

func wrapCommand(opts *flagSet, stderr io.Writer) ([]string, bool) {
	command := opts.rest()
	if len(command) == 0 {
		// A wrap without a child command would create an empty provenance shell.
		fmt.Fprintln(stderr, "wrap requires a command")
		return nil, false
	}
	return command, true
}
