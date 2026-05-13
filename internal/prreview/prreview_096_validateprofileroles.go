package prreview

func validateProfileRoles(roles []ReviewRole) (map[string]bool, error) {
	// validateProfileRoles keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	rolePlanes := map[string]bool{}
	for _, role := range roles {
		if err := validateProfileRole(role); err != nil {
			return nil, err
		}
		rolePlanes[role.Plane] = true
	}
	return rolePlanes, nil
}
