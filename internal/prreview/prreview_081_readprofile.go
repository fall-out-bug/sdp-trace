package prreview

func ReadProfile(path string) (ReviewProfile, error) {
	// ReadProfile keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	var profile ReviewProfile
	if err := readJSON(path, &profile); err != nil {
		return profile, err
	}

	return profile, validateProfile(profile)
}
