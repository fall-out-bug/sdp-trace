package harnessobs

func loadObservationSource(profilePath, sourcePath string) (Profile, []Event, string, error) {
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
