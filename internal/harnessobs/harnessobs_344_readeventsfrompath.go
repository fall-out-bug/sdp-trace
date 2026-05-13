package harnessobs

func readEventsFromPath(profilePath, sourcePath string) ([]Event, string, error) {
	// readEventsFromPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	profile, err := LoadProfile(profilePath)
	if err != nil {
		return nil, "", err
	}
	return readEvents(profile, sourcePath)
}
