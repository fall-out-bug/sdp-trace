package releaseproof

import "strings"

func evaluateManifestState(repoRoot string, manifest Manifest) verificationState {
	// Source-commit status is the root local proof; artifact checks depend on
	// that immutable source boundary being available.
	manifestSourceCommit := strings.TrimSpace(manifest.SourceCommit)
	commitStatus, reason := sourceCommitState(repoRoot, manifestSourceCommit)
	state := initialReleaseState(commitStatus)
	// Artifact results and dirty-check state can only lower confidence.
	counts, issues, artifactStatus, artifactReason := artifactVerificationState(repoRoot, manifestSourceCommit, manifest.Artifacts, state)
	state, reason = combineState(state, reason, artifactStatus, artifactReason)
	// Dirty checkout evidence is local structural evidence, not source proof.
	state, commitStatus, reason = applyDirtyState(repoRoot, state, commitStatus, reason)
	// Keep the rendered verification fields separate from decision logic.
	return verificationState{
		sourceCommit:   manifestSourceCommit,
		state:          state,
		commitStatus:   commitStatus,
		artifactStatus: artifactStatus,
		artifactCounts: counts,
		artifactIssues: issues,
		sourceReason:   reason,
	}
}

func initialReleaseState(commitStatus string) string {
	// A missing source commit breaks the immutable source boundary before any
	// artifact digest can be trusted as release proof.
	if commitStatus == StatusMissing {
		return StateCannotVerify
	}
	return StatePass
}
