package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/repoobserver"
)

func repoObserverOptionsFromFlags(opts *flagSet) repoobserver.Options {
	// Convert flags directly into repoobserver options so install/doctor keep a
	// single source of truth for profile semantics.
	return repoobserver.Options{
		Profile:      opts.stringValue("profile"),
		RepositoryID: opts.stringValue("repository-id"),
		Write:        opts.boolValue("write"),
		Force:        opts.boolValue("force"),
	}
}
