package packet

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// Authority metadata records who wrote an entry and what source ref it is bound
// to. Source refs are trimmed so whitespace cannot masquerade as evidence.
func authorityEntry(entry BundleEntry, actor, writeAuthority, generatedBy, sourceCommitState, sourceRef string) BundleEntry {
	entry.Actor = actor
	entry.WriteAuthority = writeAuthority
	entry.GeneratedBy = generatedBy
	entry.SourceCommitState = sourceCommitState
	entry.SourceRef = strings.TrimSpace(sourceRef)
	return entry
}

// Bundle entries default to present, not_needed redaction and a deterministic
// digest placeholder. Resolver redaction happens before the placeholder digest
// is calculated so secrets do not influence checked-in evidence text.
func bundleEntry(ref, sourceClass, resolver, retainedForm string) BundleEntry {
	resolver = redactSecretLike(resolver)

	return BundleEntry{
		Ref:             ref,
		SourceClass:     sourceClass,
		Digest:          digestPlaceholder(ref + resolver),
		RetainedForm:    retainedForm,
		RedactionStatus: "not_needed",
		Resolver:        resolver,
		ArtifactAccess:  "present",
	}
}

// The redaction check is deliberately small and conservative; it prevents
// obvious secret-like resolver strings from being persisted in generated
// examples or packets.
func redactSecretLike(value string) string {
	redacted := value
	for _, marker := range []string{"SECRET", "TOKEN", "Authorization:"} {
		if strings.Contains(strings.ToUpper(redacted), strings.ToUpper(marker)) {
			return "[redacted-secret]"
		}
	}
	return redacted
}

func digestPlaceholder(value string) string {
	sum := sha256.Sum256([]byte(value))
	return "sha256:" + fmt.Sprintf("%x", sum[:])
}
