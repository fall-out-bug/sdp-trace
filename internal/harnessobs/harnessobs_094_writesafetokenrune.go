package harnessobs

import (
	"strings"
)

func writeSafeTokenRune(b *strings.Builder, r rune) {
	// writeSafeTokenRune keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if safeTokenRune(r) {
		b.WriteRune(r)
		return
	}

	b.WriteByte('-')
}
