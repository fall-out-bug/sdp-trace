package main

import (
	"fmt"
	"io"
)

func requireInstallRepoObserverFlags(opts *flagSet, stderr io.Writer) (*flagSet, int, bool) {
	if len(opts.rest()) != 0 {
		// Repo-observer install is fully flag-driven; no positional repository
		// path is interpreted here.
		fmt.Fprintln(stderr, "install repo-observer accepts only flags")
		return nil, exitUsage, false
	}
	return opts, 0, true
}
