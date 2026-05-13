package releaseproof

import (
	"os/exec"
	"strings"
)

func sourceCommitState(repoRoot, sourceCommit string) (string, string) {
	// The source commit is the anchor for every source-bound artifact check; if
	// it cannot resolve locally, artifact claims stay unverified.
	if sourceCommit == "" {
		return StatusMissing, "manifest source_commit is missing"
	}
	if !sourceCommitExists(repoRoot, sourceCommit) {
		return StatusMissing, "manifest source_commit could not be resolved from git"
	}
	return StatusMatched, "source commit contains every manifest artifact path with matching digest"
}

func combineState(state, reason, artifactStatus, artifactReason string) (string, string) {
	// Artifact verification can only lower confidence; it never upgrades a
	// missing source commit or previously unverified release proof.
	if state == StateCannotVerify || artifactStatus == StatusMatched {
		return state, reason
	}
	return StateFail, artifactReason
}

func applyDirtyState(repoRoot, state, commitStatus, reason string) (string, string, string) {
	// A dirty checkout is local structural evidence only, so it blocks a source
	// match without turning external trust green.
	if state == StateCannotVerify || !worktreeDirty(repoRoot) {
		return state, commitStatus, reason
	}
	return StateFail, StatusMismatch, "dirty checkout cannot support source-bound local release proof"
}

func sourceCommitExists(repoRoot, sourceCommit string) bool {
	// Git object resolution is the immutable-source boundary for this local
	// verifier; a missing object keeps the release verdict at cannot_verify.
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
