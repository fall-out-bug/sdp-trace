package harnessobs

func LoadProfile(path string) (Profile, error) {
	// Harness profiles remain ordinary local JSON inputs until validation
	// succeeds; loading them does not create proof or authority.
	var profile Profile

	if err := readExistingJSON(path, &profile); err != nil {
		return Profile{}, err
	}
	if err := validateProfile(profile); err != nil {
		return Profile{}, err
	}
	return profile, nil
}

func LoadSessionProfile(path string) (SessionProfile, error) {
	// Session profiles use strict JSON because unknown or trailing fields can
	// change what raw events, setup actions, or isolation rules are trusted.
	var profile SessionProfile

	if err := readExistingJSONStrict(path, &profile); err != nil {
		return SessionProfile{}, err
	}
	if err := validateSessionProfile(&profile); err != nil {
		return SessionProfile{}, err
	}
	return profile, nil
}
