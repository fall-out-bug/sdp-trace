package harnessobs

func (scanner *shellFieldScanner) consumeEscaped(r rune) bool {
	// consumeEscaped keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !scanner.escaped {
		return false
	}
	if r != '\n' {

		scanner.current.WriteRune('\\')
		scanner.current.WriteRune(r)
	}
	scanner.escaped = false
	return true
}
