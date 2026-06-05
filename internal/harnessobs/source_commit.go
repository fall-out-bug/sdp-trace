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

// currentSourceCommitState converts local checkout discovery into explicit
// evidence state without treating missing git context as success.
func currentSourceCommitState() (string, string) {
	commit := sourceCommit()
	if commit == "" {
		return "", StateCannotVerify
	}
	return commit, StatePass
}
