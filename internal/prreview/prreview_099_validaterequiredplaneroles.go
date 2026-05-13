package prreview

import (
	"fmt"
)

func validateRequiredPlaneRoles(requiredPlanes []string, rolePlanes map[string]bool) error {
	// validateRequiredPlaneRoles keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	for _, plane := range requiredPlanes {
		if !rolePlanes[plane] {
			return fmt.Errorf("profile_required_plane_without_role: %s", plane)
		}
	}
	return nil
}
