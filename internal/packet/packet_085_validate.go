package packet

func (v *bundleValidator) validate() Validation {
	// validate keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	v.validateMetadata()
	v.indexManifest()
	v.validateRows()
	v.validateFindingsAndGaps()
	state := StatePass
	if len(v.errors) > 0 {
		state = StateFail
	}
	return Validation{State: state, Errors: v.errors}
}
