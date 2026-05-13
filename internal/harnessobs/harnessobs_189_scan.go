package harnessobs

func (scanner *shellFieldScanner) scan(r rune) {
	// scan keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if scanner.consumeEscaped(r) {
		return
	}
	if scanner.startsEscape(r) {

		scanner.escaped = true
		return
	}
	if scanner.consumeQuoted(r) {
		return
	}

	scanner.consumeUnquoted(r)
}
