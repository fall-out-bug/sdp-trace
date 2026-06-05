package prreview

import (
	"errors"
	"fmt"
)

// validateProfileRoles validates each declared role and returns the planes
// covered by executable review roles.
func validateProfileRoles(roles []ReviewRole) (map[string]bool, error) {
	rolePlanes := map[string]bool{}
	for _, role := range roles {
		if err := validateProfileRole(role); err != nil {
			return nil, err
		}
		rolePlanes[role.Plane] = true
	}
	return rolePlanes, nil
}

// validateProfileRole keeps role execution portable by accepting only complete
// roles with supported runner names.
func validateProfileRole(role ReviewRole) error {
	if profileRoleMissingRequiredField(role) {
		return errors.New("profile_role_requires_id_plane_runner")
	}
	if !validRunner(role.Runner) {
		return fmt.Errorf("profile_role_invalid_runner: %s", role.Runner)
	}
	return nil
}

// profileRoleMissingRequiredField detects incomplete role declarations before
// runner selection can infer behavior from missing data.
func profileRoleMissingRequiredField(role ReviewRole) bool {
	return role.RoleID == "" || role.Plane == "" || role.Runner == ""
}

// validateRequiredPlaneRoles verifies the profile can execute every required
// review plane.
func validateRequiredPlaneRoles(requiredPlanes []string, rolePlanes map[string]bool) error {
	for _, plane := range requiredPlanes {
		if !rolePlanes[plane] {
			return fmt.Errorf("profile_required_plane_without_role: %s", plane)
		}
	}
	return nil
}
