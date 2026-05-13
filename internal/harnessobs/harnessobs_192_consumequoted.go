package harnessobs

func (scanner *shellFieldScanner) consumeQuoted(r rune) bool {
	// consumeQuoted keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if scanner.quote == 0 {
		return false
	}
	if r == scanner.quote {

		scanner.quote = 0
		return true
	}
	scanner.current.WriteRune(r)
	return true
}
