package harnessobs

import (
	"os/exec"
	"strings"
)

// sourceCommit records the current checkout HEAD and fails closed outside git.
func sourceCommit() string {
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	data, err := cmd.Output()
	if err != nil {
		return ""
	}

	commit := strings.TrimSpace(string(data))
	if sourceCommitHash(commit) {
		return commit
	}
	return ""
}
