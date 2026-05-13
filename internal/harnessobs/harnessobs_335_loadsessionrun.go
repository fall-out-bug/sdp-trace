package harnessobs

func LoadSessionRun(path string) (SessionRun, error) {
	// LoadSessionRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var run SessionRun

	if err := readExistingJSON(path, &run); err != nil {
		return SessionRun{}, err
	}
	if err := validateLoadedSessionRun(run); err != nil {
		return SessionRun{}, err
	}
	return run, nil
}
