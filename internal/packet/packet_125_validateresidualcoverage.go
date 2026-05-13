package packet

func (v *bundleValidator) validateResidualCoverage(rows map[string]Row) {
	// validateResidualCoverage keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, row := range rows {
		if residualCoverageExempt(row) {
			continue
		}

		if !gapForRow(v.bundle.Packet.ResidualGaps, row.ID) {
			v.add("%s non-pass row requires residual gap explanation", row.ID)
		}
	}
}
