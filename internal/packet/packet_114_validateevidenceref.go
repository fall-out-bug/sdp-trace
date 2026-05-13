package packet

import (
	"strings"
)

func (v *bundleValidator) validateEvidenceRef(rowID, state, ref string) {
	// validateEvidenceRef keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	entry, ok := v.entryByRef[ref]
	if !ok {

		v.add("%s evidence ref %q is absent from manifest", rowID, ref)
		return
	}
	if strings.TrimSpace(v.resolverByRef[ref]) == "" {

		v.add("%s evidence ref %q has no resolver entry", rowID, ref)
	}
	if state != StatePass {
		return
	}
	v.validatePassEvidenceRef(rowID, ref, entry)
}
