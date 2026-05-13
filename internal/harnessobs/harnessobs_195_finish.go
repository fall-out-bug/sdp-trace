package harnessobs

func (scanner *shellFieldScanner) finish() []string {
	// finish keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if scanner.escaped {

		scanner.current.WriteRune('\\')
	}
	scanner.flush()
	return scanner.fields
}
