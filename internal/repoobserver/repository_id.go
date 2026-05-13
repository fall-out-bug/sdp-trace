package repoobserver

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func derivedRepositoryID(root string) string {
	// Derived IDs hash sanitized remote data to avoid rendering credentials.
	origin := gitOutput(root, "config", "--get", "remote.origin.url")
	if strings.TrimSpace(origin) == "" {
		origin = "current_repository"
	}
	sanitized := sanitizeOrigin(origin)
	sum := sha256.Sum256([]byte(sanitized))
	return "repo_" + hex.EncodeToString(sum[:])[:16]
}
