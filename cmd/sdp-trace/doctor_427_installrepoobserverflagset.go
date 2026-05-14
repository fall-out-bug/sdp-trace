package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/repoobserver"
)

func installRepoObserverFlagSet() *flagSet {
	opts := &flagSet{name: "install repo-observer"}
	// Default to the only supported portable repo-observer profile.
	opts.setString("profile", repoobserver.ProfileGithubActionsGitHooksV1)
	opts.setString("repository-id", "")
	opts.setString("out", "")
	opts.setBool("write", false)
	opts.setBool("force", false)
	return opts
}
