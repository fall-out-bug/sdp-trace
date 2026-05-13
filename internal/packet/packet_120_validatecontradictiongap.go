package packet

func (v *bundleValidator) validateContradictionGap(rowID string) {
	// validateContradictionGap keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if !gapForRow(v.bundle.Packet.ResidualGaps, rowID) {

		v.add("%s contradictory evidence requires residual gap explanation", rowID)
	}
}
