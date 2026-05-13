package witness

import (
	"os"
	"strings"
)

func ambientCIEnvPresent(kind string) bool {
	// Ambient variables explain why an envelope is required, but this helper
	// never upgrades trust by itself.
	prefixes := map[string][]string{
		KindGitLabCI:  {"GITLAB_CI", "CI_PIPELINE_ID", "CI_JOB_ID", "CI_COMMIT_SHA"},
		KindBuildkite: {"BUILDKITE", "BUILDKITE_BUILD_ID", "BUILDKITE_JOB_ID", "BUILDKITE_COMMIT"},
	}
	for _, key := range prefixes[kind] {
		if strings.TrimSpace(os.Getenv(key)) != "" {
			return true
		}
	}
	return false
}
