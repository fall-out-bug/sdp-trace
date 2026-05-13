package prreview

import (
	"crypto/sha256"
	"encoding/hex"

	"os/exec"

	"strings"
)

func captureWorkingTreeBaseline(workDir string) (*workingTreeBaseline, error) {
	// captureWorkingTreeBaseline keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

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
