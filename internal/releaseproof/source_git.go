package releaseproof

import (
	"os/exec"
	"strings"
)

func sourceCommitExists(repoRoot, sourceCommit string) bool {
	// Git object resolution is the immutable-source boundary for this local
	// verifier; a missing object keeps the release verdict at cannot_verify.
	// sourceCommit is validated as a 40-char lowercase hex SHA before this call.
	// #nosec G204
	cmd := exec.Command("git", "cat-file", "-e", sourceCommit+"^{commit}")
	cmd.Dir = repoRoot
	return cmd.Run() == nil
}

func worktreeDirty(repoRoot string) bool {
	// Treat git status failures as dirty so command failures cannot accidentally
	// promote a source-bound proof.
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		return true
	}
	return strings.TrimSpace(string(out)) != ""
}
