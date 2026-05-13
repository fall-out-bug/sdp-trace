package harnessobs

func LoadSessionProfile(path string) (SessionProfile, error) {
	// LoadSessionProfile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var profile SessionProfile

	if err := readExistingJSONStrict(path, &profile); err != nil {
		return SessionProfile{}, err
	}
	if err := validateSessionProfile(&profile); err != nil {
		return SessionProfile{}, err
	}
	return profile, nil
}
