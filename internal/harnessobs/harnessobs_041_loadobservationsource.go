package harnessobs

func loadObservationSource(profilePath, sourcePath string) (Profile, []Event, string, error) {
	// loadObservationSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profile, err := LoadProfile(profilePath)
	if err != nil {
		return Profile{}, nil, "", err
	}

	events, sourceDigest, err := readEvents(profile, sourcePath)
	if err != nil {
		return Profile{}, nil, "", err
	}
	return profile, events, sourceDigest, nil
}
