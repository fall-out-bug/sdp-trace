package packet

func (v *bundleValidator) validatePassRowEvidence(row Row) {
	// validatePassRowEvidence keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if row.State == StatePass && len(row.EvidenceRefs) == 0 {

		v.add("%s pass requires retained evidence refs", row.ID)
	}
}
