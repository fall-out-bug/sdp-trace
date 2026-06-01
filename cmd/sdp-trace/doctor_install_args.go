package main

import (
	"fmt"
	"io"
)

func parseInstallRepoObserverArgs(args []string, stdout, stderr io.Writer) (*flagSet, int, bool) {
	if isHelp(args) {
		printUsage(stdout)
		return nil, 0, false
	}
	if !hasInstallRepoObserverSubcommand(args) {
		// The installer namespace is intentionally closed until another portable
		// installer contract is added.
		// Keep install scoped to repo-observer so future installers do not share
		// ambiguous flag contracts.
		fmt.Fprintln(stderr, "install requires repo-observer")
		return nil, exitUsage, false
	}
	opts := installRepoObserverFlagSet()
	if err := opts.parse(args[1:]); err != nil {
		// Parse only the arguments after the required repo-observer verb.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return requireInstallRepoObserverFlags(opts, stderr)
}

func requireInstallRepoObserverFlags(opts *flagSet, stderr io.Writer) (*flagSet, int, bool) {
	if len(opts.rest()) != 0 {
		// Repo-observer install is fully flag-driven; no positional repository
		// path is interpreted here.
		fmt.Fprintln(stderr, "install repo-observer accepts only flags")
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func hasInstallRepoObserverSubcommand(args []string) bool {
	return len(args) != 0 && args[0] == "repo-observer"
}
