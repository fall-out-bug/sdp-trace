package packet

func (v *bundleValidator) requireRowsPresent(rows map[string]Row) {
	// requireRowsPresent keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, id := range RequiredRows {
		if rows[id].ID == "" {

			v.add("missing required row %q", id)
		}
	}
}
