package packet

func (v *bundleValidator) validateRowState(row Row) {
	// validateRowState keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if !states[row.State] {

		v.add("%s has unknown state %q", row.ID, row.State)
	}
}
