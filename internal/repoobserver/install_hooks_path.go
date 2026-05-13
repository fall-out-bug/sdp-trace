package repoobserver

import (
	"fmt"
	"strings"
)

func appendHooksPathSummary(opts Options, summaries []DiffSummary) ([]DiffSummary, error) {
	// hooksPath changes are local git config mutations and are summarized
	// separately from repository file writes.
	// Existing non-default hook paths are summarized only in force mode after the
	// caller has accepted replacement.
	previousHooksPath := strings.TrimSpace(gitOutput(opts.RepoRoot, "config", "--get", "core.hooksPath"))
	if opts.Force && isDifferentHooksPath(previousHooksPath) {
		summaries = append(summaries, DiffSummary{
			Path:    "git_config:core.hooksPath",
			Action:  "overwrite_hooks_path",
			Before:  safeRef(previousHooksPath),
			After:   ".githooks",
			Summary: "replace local checkout hooks path reference",
		})
	}
	if err := runGit(opts.RepoRoot, "config", "core.hooksPath", ".githooks"); err != nil {
		return summaries, err
	}
	return summaries, nil
}

func isDifferentHooksPath(path string) bool {
	return path != "" && path != ".githooks"
}

func ensureNoUnsafeHooksPath(opts Options) error {
	// Existing non-.githooks values require --force so user hook configuration is
	// not silently replaced.
	value := strings.TrimSpace(gitOutput(opts.RepoRoot, "config", "--get", "core.hooksPath"))
	if value == "" || value == ".githooks" || opts.Force {
		return nil
	}
	return fmt.Errorf("%s: core.hooksPath is %s; use --force only after reviewing existing hooks", ReasonHooksPathMismatch, safeRef(value))
}
