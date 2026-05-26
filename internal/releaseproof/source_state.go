package releaseproof

import (
	"regexp"
)

var commitSHAPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)

func sourceCommitState(repoRoot, sourceCommit string) (string, string) {
	// The source commit is the anchor for every source-bound artifact check; if
	// it cannot resolve locally, artifact claims stay unverified.
	if sourceCommit == "" {
		return StatusMissing, "manifest source_commit is missing"
	}
	if !isValidCommitSHA(sourceCommit) {
		return StatusMissing, "manifest source_commit is not a valid immutable commit SHA"
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

func isValidCommitSHA(ref string) bool {
	// Accept only immutable 40-character lowercase hex commit object
	// identifiers. Reject branch names, symbolic refs, revspec suffixes,
	// pathspecs, flags, and any ref that can resolve to a non-commit object.
	return commitSHAPattern.MatchString(ref)
}
