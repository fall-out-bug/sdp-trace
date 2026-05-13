package prreview

func validateProfile(profile ReviewProfile) error {
	// validateProfile keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if err := validateProfileHeader(profile); err != nil {
		return err
	}
	rolePlanes, err := validateProfileRoles(profile.Roles)
	if err != nil {
		return err
	}
	return validateRequiredPlaneRoles(profile.RequiredPlanes, rolePlanes)
}
