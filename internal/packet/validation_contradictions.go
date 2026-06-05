package packet

// Contradiction validation walks manifest entries rather than packet rows so
// every retained contradictory artifact can force its target row back to
// partial plus a residual-gap explanation.
func (v *bundleValidator) validateContradictions(rows map[string]Row) {
	for _, entry := range v.entryByRef {
		v.validateContradiction(rows, entry)
	}
}

// Explicit contradiction row IDs are authoritative; otherwise attribution falls
// back to exact evidence-ref matching through rowIDForRef.
func (v *bundleValidator) validateContradiction(rows map[string]Row, entry BundleEntry) {
	rowID := contradictionRowID(rows, entry)
	if !hasContradictionTarget(entry, rowID) {
		return
	}
	row := rows[rowID]
	v.validateContradictionState(rowID, row)
	v.validateContradictionGap(rowID)
}

func hasContradictionTarget(entry BundleEntry, rowID string) bool {
	return entry.ContradictsRef != "" && rowID != ""
}

// A manifest entry may name the row directly to avoid ambiguous ref ownership;
// absent that, the row lookup preserves required-row precedence.
func contradictionRowID(rows map[string]Row, entry BundleEntry) string {
	if entry.ContradictsRowID != "" {
		return entry.ContradictsRowID
	}
	return rowIDForRef(rows, entry.Ref)
}
