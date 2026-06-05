package packet

import (
	"strings"
	"time"
)

var unverifiableArtifactAccess = map[string]bool{
	"expired":         true,
	"inaccessible":    true,
	"malformed":       true,
	"not_assessed":    true,
	StateCannotVerify: true,
}

// Demo-first gates accept only evidence that can be replayed from a resolver,
// has a digest, has not expired, and is still retained or otherwise accessible.
func demoUsableEntry(entry BundleEntry, now time.Time) bool {
	return entryHasResolverAndDigest(entry) && !entryExpired(entry, now) && !passRefUnverifiable(entry)
}

// A retained evidence entry needs both addressability and integrity material.
func entryHasResolverAndDigest(entry BundleEntry) bool {
	return strings.TrimSpace(entry.Resolver) != "" && strings.TrimSpace(entry.Digest) != ""
}

// Synthetic digests are placeholders for generated examples, not proof.
func syntheticEntryDigest(entry BundleEntry) bool {
	return strings.TrimSpace(entry.Digest) == "" || entry.Digest == digestPlaceholder(entry.Ref+entry.Resolver)
}

// Missing expiry means no expiry claim; malformed expiry is treated as unusable
// because a verifier cannot safely establish freshness.
func entryExpired(entry BundleEntry, now time.Time) bool {
	if strings.TrimSpace(entry.ExpiresAt) == "" {
		return false
	}
	expiresAt, err := time.Parse(time.RFC3339, entry.ExpiresAt)
	if err != nil {
		return true
	}
	return !expiresAt.After(now)
}

// Passing evidence is unverifiable when retention or access says the artifact
// cannot be replayed by a reviewer.
func passRefUnverifiable(entry BundleEntry) bool {
	if entry.RedactionStatus == StateCannotVerify || entry.RetainedForm == "not_retained" {
		return true
	}
	return artifactAccessUnverifiable(entry.ArtifactAccess)
}

// Unknown access values remain forward-compatible and do not fail pass rows.
func artifactAccessUnverifiable(access string) bool {
	return unverifiableArtifactAccess[access]
}
