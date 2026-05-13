package packet

func (v *bundleValidator) validateRowEvidenceRefs(row Row) {
	// validateRowEvidenceRefs keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, ref := range row.EvidenceRefs {

		v.validateEvidenceRef(row.ID, row.State, ref)
	}
}
