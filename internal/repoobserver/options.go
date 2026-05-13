package repoobserver

import (
	"fmt"
	"strings"
	"time"
)

const ProfileGithubActionsGitHooksV1 = "github-actions-git-hooks-v1"

type Options struct {
	RepoRoot     string
	Profile      string
	RepositoryID string
	Write        bool
	Force        bool
	Now          time.Time
}

func normalizeOptions(opts Options) (Options, error) {
	// Defaults, absolute paths, and repository identity all pass through one
	// normalization path before any observation or install work.
	opts = withDefaultProfile(opts)
	if err := validateProfile(opts.Profile); err != nil {
		return Options{}, err
	}
	opts, err := withAbsoluteRepoRoot(opts)
	if err != nil {
		return Options{}, err
	}
	opts = withDefaultNow(opts)
	return opts, validateRepositoryID(opts.RepositoryID, "repository id must match [A-Za-z0-9_.-]+")
}

func withDefaultProfile(opts Options) Options {
	if strings.TrimSpace(opts.Profile) == "" {
		// Empty profile selects the only portable repo-observer contract rather
		// than leaving behavior harness-specific.
		opts.Profile = ProfileGithubActionsGitHooksV1
	}
	return opts
}

func validateProfile(profile string) error {
	if profile != ProfileGithubActionsGitHooksV1 {
		// Unknown profiles would imply setup semantics this portable tool cannot
		// verify.
		return fmt.Errorf("repo observer requires --profile %s", ProfileGithubActionsGitHooksV1)
	}
	return nil
}
