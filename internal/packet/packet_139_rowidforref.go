package packet

func rowIDForRef(rows map[string]Row, ref string) string {
	// rowIDForRef keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if id := requiredRowIDForRef(rows, ref); id != "" {
		return id
	}
	return extensionRowIDForRef(rows, ref)
}
