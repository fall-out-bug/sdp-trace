package harnessobs

func (scanner *shellFieldScanner) flush() {
	// flush keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if scanner.current.Len() == 0 {
		return
	}

	scanner.fields = append(scanner.fields, scanner.current.String())
	scanner.current.Reset()
}
