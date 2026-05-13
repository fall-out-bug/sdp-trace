package harnessobs

func ensureLineFileRule(path, line string) error {
	// ensureLineFileRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	lines, err := readOptionalLines(path)
	if err != nil {
		return err
	}

	for _, existing := range lines {
		if existing == line {
			return nil
		}
	}
	lines = append(lines, line)
	return writeLines(path, lines)
}
