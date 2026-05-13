package packet

func (v *bundleValidator) validateRows() {
	// validateRows keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	rows := map[string]Row{}
	v.indexRows(rows)
	v.requireRowsPresent(rows)
	v.rows = rows

	v.validateContradictions(rows)
	v.validateResidualCoverage(rows)
}
