package packet

func (v *bundleValidator) validateRowID(rowID string, rows map[string]Row) bool {
	// validateRowID keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	if !requiredRow(rowID) {
		v.add("unknown row id %q", rowID)
		return false
	}
	if rows[rowID].ID != "" {
		v.add("duplicate row id %q", rowID)
	}
	return true
}
