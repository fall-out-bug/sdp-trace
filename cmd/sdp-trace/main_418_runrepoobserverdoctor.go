package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/repoobserver"
)

func runRepoObserverDoctor(opts *flagSet, stdout, stderr io.Writer) int {
	if opts.stringValue("profile") != repoobserver.ProfileGithubActionsGitHooksV1 {
		// The CLI exposes only the portable GitHub Actions/git-hooks profile.
		fmt.Fprintf(stderr, "doctor --profile requires %s\n", repoobserver.ProfileGithubActionsGitHooksV1)
		return exitUsage
	}
	// Doctor is read-only; it reports install/proof state without modifying the
	// repository.
	status, err := repoobserver.Doctor(repoobserver.Options{Profile: opts.stringValue("profile")})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return writeRepoObserverDoctor(opts, status, stdout, stderr)
}
