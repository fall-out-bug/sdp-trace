package main

import (
	"fmt"
	"io"
)

func requireRest(opts *flagSet, stderr io.Writer, message string) bool {
	if len(opts.rest()) != 0 {
		return true
	}
	// Commands after `--` are part of the replay boundary; missing rest args
	// would record feedback without a target command.
	fmt.Fprintln(stderr, message)
	return false
}
