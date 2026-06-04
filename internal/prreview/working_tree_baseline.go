package prreview

import (
	"crypto/sha256"
	"encoding/hex"
	"os/exec"
	"strings"
)

type workingTreeBaseline struct {
	Count  int
	Digest string
}

// captureWorkingTreeBaseline hashes porcelain status so OpenCode mutation
// checks can compare pre-run and post-run worktree state.
func captureWorkingTreeBaseline(workDir string) (*workingTreeBaseline, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = workDir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	count := 0
	if strings.TrimSpace(string(output)) != "" {
		count = len(lines)
	}
	sum := sha256.Sum256(output)
	return &workingTreeBaseline{Count: count, Digest: hex.EncodeToString(sum[:])}, nil
}
