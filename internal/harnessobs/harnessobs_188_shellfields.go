package harnessobs

func shellFields(command string) []string {
	// shellFields keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	scanner := shellFieldScanner{}
	for _, r := range command {
		scanner.scan(r)
	}
	return scanner.finish()
}
