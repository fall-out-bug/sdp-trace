package packet

func (v *bundleValidator) validateContradiction(rows map[string]Row, entry BundleEntry) {
	// validateContradiction keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	rowID := contradictionRowID(rows, entry)
	if !hasContradictionTarget(entry, rowID) {
		return
	}
	row := rows[rowID]
	v.validateContradictionState(rowID, row)
	v.validateContradictionGap(rowID)
}
