package packet

import (
	"strings"
)

func authorityEntry(entry BundleEntry, actor, writeAuthority, generatedBy, sourceCommitState, sourceRef string) BundleEntry {
	// authorityEntry keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	entry.Actor = actor
	entry.WriteAuthority = writeAuthority
	entry.GeneratedBy = generatedBy
	entry.SourceCommitState = sourceCommitState
	entry.SourceRef = strings.TrimSpace(sourceRef)
	return entry
}
