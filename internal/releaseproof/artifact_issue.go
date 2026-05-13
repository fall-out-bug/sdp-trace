package releaseproof

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
)

func artifactIssue(repoRoot, sourceCommit string, artifact ManifestArtifact) (ArtifactIssue, bool) {
	// Clean the path used for git object lookup, but keep the original
	// manifest path in issues so reports match the signed obligation text.
	path := filepath.Clean(artifact.Path)
	data, err := artifactBytes(repoRoot, sourceCommit, path)
	if err != nil {
		// Missing source-commit bytes are stronger than digest mismatch because
		// there is no artifact content to compare.
		return ArtifactIssue{Path: artifact.Path, Issue: StatusMissing, Expected: artifact.SHA256}, true
	}
	sum := sha256.Sum256(data)
	actual := hex.EncodeToString(sum[:])
	// Hex case is formatting, not proof content; compare digest values
	// case-insensitively while reporting the canonical lowercase actual.
	if strings.EqualFold(actual, artifact.SHA256) {
		return ArtifactIssue{}, false
	}
	return ArtifactIssue{Path: artifact.Path, Issue: StatusMismatch, Expected: artifact.SHA256, Actual: actual}, true
}
