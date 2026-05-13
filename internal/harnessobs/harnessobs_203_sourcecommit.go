package harnessobs

import (
	"os/exec"

	"regexp"
	"strings"
)

func sourceCommit() string {
	// sourceCommit keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	data, err := cmd.Output()
	if err != nil {
		return ""
	}

	commit := strings.TrimSpace(string(data))
	if regexp.MustCompile(`^[a-f0-9]{40}$`).MatchString(commit) {

		return commit
	}
	return ""
}
