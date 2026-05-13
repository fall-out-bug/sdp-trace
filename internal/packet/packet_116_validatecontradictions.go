package packet

func (v *bundleValidator) validateContradictions(rows map[string]Row) {
	// validateContradictions keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, entry := range v.entryByRef {

		v.validateContradiction(rows, entry)
	}
}
