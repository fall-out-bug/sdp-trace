package harnessobs

func (scanner *shellFieldScanner) consumeUnquoted(r rune) {
	// consumeUnquoted keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch {
	case shellQuote(r):

		scanner.quote = r
	case shellFieldSeparator(r):
		scanner.flush()
	default:
		scanner.current.WriteRune(r)
	}
}
