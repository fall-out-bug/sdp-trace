package repoobserver

import (
	"errors"
	"fmt"
	"os"
	"regexp"
)

const ReasonUnsafeOutputRefused = "unsafe_output_refused"

var safeIDPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

func withConfiguredRepositoryID(opts Options) (Options, error) {
	// Repository identity is fixed before status construction so every surface
	// is rendered against the same source label.
	if opts.RepositoryID != "" {
		// Explicit repository identity wins over local config so callers can
		// bind observations to the intended source.
		return opts, nil
	}
	return withConfigFileRepositoryID(opts)
}

func withConfigFileRepositoryID(opts Options) (Options, error) {
	// Missing config is allowed for first install/doctor runs; the caller can
	// still use a derived repository ID.
	config, err := LoadConfig(opts.RepoRoot)
	if errors.Is(err, os.ErrNotExist) {
		return opts, nil
	}
	if err != nil {
		return Options{}, err
	}
	if config.RepositoryID != "" {
		opts.RepositoryID = config.RepositoryID
	}
	return opts, nil
}

func validateConfig(config Config) (Config, error) {
	if config.Profile != "" && config.Profile != ProfileGithubActionsGitHooksV1 {
		// Config files cannot opt into unimplemented observer profiles.
		return Config{}, fmt.Errorf("%s: unsupported repo observer profile in .sdp-trace/config.json", ReasonUnsafeOutputRefused)
	}
	return config, validateRepositoryID(config.RepositoryID, "repository id in .sdp-trace/config.json must match [A-Za-z0-9_.-]+")
}

func validateRepositoryID(repositoryID, message string) error {
	if repositoryID != "" && !safeIDPattern.MatchString(repositoryID) {
		// Repository IDs are rendered in reports and config, so reject unsafe
		// labels at the boundary.
		return fmt.Errorf("%s: %s", ReasonUnsafeOutputRefused, message)
	}
	return nil
}
