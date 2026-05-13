package packet

func (v *bundleValidator) requireKnown(known map[string]bool, value string, format string) {
	// requireKnown keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if !known[value] {

		v.add(format, value)
	}
}
