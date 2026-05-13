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
