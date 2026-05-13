package harnessobs

func validUnavailableField(field UnavailableField) bool {
	// validUnavailableField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return safeIDPattern.MatchString(field.Field) &&
		field.State == StateNotAssessed &&
		safeIDPattern.MatchString(field.ReasonCode)
}
