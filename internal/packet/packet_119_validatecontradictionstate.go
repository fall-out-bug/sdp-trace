package packet

func (v *bundleValidator) validateContradictionState(rowID string, row Row) {
	// validateContradictionState keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if row.State != StatePartial {

		v.add("%s has contradictory evidence but state is %s, want partial", rowID, row.State)
	}
}
