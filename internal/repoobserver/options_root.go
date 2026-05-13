package repoobserver

import (
	"path/filepath"
	"strings"
)

func withAbsoluteRepoRoot(opts Options) (Options, error) {
	// Blank roots are resolved through git and then made absolute for later
	// containment checks.
	// The absolute path is local structural evidence only and is rendered as an
	// abstract repository root in user-facing status.
	if strings.TrimSpace(opts.RepoRoot) == "" {
		root, err := repoRoot(".")
		if err != nil {
			return Options{}, err
		}
		opts.RepoRoot = root
	}
	abs, err := filepath.Abs(opts.RepoRoot)
	if err != nil {
		return Options{}, err
	}
	opts.RepoRoot = abs
	return opts, nil
}
