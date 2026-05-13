package harnessobs

func LoadProfile(path string) (Profile, error) {
	// LoadProfile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var profile Profile

	if err := readExistingJSON(path, &profile); err != nil {
		return Profile{}, err
	}
	if err := validateProfile(profile); err != nil {
		return Profile{}, err
	}
	return profile, nil
}
