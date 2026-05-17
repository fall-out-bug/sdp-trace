package releaseproof

import "os/exec"

type ArtifactCounts struct {
	Checked    int `json:"checked"`
	Missing    int `json:"missing"`
	Mismatched int `json:"mismatched"`
}

type ArtifactIssue struct {
	Path     string `json:"path"`
	Issue    string `json:"issue"`
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
}

func artifactVerificationState(repoRoot, sourceCommit string, artifacts []ManifestArtifact, state string) (ArtifactCounts, []ArtifactIssue, string, string) {
	// Do not inspect artifacts when the source anchor is absent; that would turn
	// missing source proof into misleading artifact evidence.
	if state == StateCannotVerify {
		return ArtifactCounts{}, nil, StatusNotAssessed, "manifest artifacts were not checked because source_commit cannot be verified"
	}
	counts, issues := artifactCounts(repoRoot, sourceCommit, artifacts)
	artifactStatus, artifactReason := artifactState(counts)
	return counts, issues, artifactStatus, artifactReason
}

func artifactState(counts ArtifactCounts) (string, string) {
	// Missing artifact paths are reported before digest mismatches because the
	// verifier cannot compare bytes that were not present in the source commit.
	if counts.Missing > 0 {
		return StatusMissing, "manifest artifact paths are missing from the current source checkout"
	}
	if counts.Mismatched > 0 {
		return StatusMismatch, "manifest artifact digests do not match the current source checkout"
	}
	return StatusMatched, ""
}

func artifactBytes(repoRoot, sourceCommit, path string) ([]byte, error) {
	// Read artifacts from the manifest source commit, not the dirty checkout,
	// so local edits cannot satisfy source-bound release proof.
	if !isValidCommitSHA(sourceCommit) {
		return nil, exec.ErrNotFound
	}
	// sourceCommit is validated as a 40-char lowercase hex SHA; path is cleaned
	// by filepath.Clean before this call.
	// #nosec G204
	cmd := exec.Command("git", "show", sourceCommit+":"+path)
	cmd.Dir = repoRoot
	return cmd.Output()
}
