package packet

func (v *bundleValidator) validateRow(row Row) {
	// validateRow keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	v.validateRowRequiredFields(row)
	v.validateRowReason(row)
	v.validatePassRowEvidence(row)
	v.validateRowEvidenceRefs(row)
}
