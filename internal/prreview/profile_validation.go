package prreview

import (
	"errors"
	"fmt"
	"strings"
)

// validateProfile checks the profile header and verifies every required plane
// has at least one executable review role.
func validateProfile(profile ReviewProfile) error {
	if err := validateProfileHeader(profile); err != nil {
		return err
	}
	rolePlanes, err := validateProfileRoles(profile.Roles)
	if err != nil {
		return err
	}
	return validateRequiredPlaneRoles(profile.RequiredPlanes, rolePlanes)
}

// validateProfileHeader enforces the profile schema contract while preserving
// the legacy empty schema version accepted by existing fixtures.
func validateProfileHeader(profile ReviewProfile) error {
	if profile.SchemaVersion != "" && profile.SchemaVersion != SchemaVersionProfile {
		return fmt.Errorf("invalid_profile_schema_version: %s", profile.SchemaVersion)
	}
	return requireProfileFields(profile)
}

// requireProfileFields checks the minimum profile shape needed to plan review
// execution deterministically.
func requireProfileFields(profile ReviewProfile) error {
	if strings.TrimSpace(profile.ProfileID) == "" {
		return errors.New("profile_requires_profile_id")
	}
	if len(profile.RequiredPlanes) == 0 {
		return errors.New("profile_requires_required_planes")
	}
	if len(profile.Roles) == 0 {
		return errors.New("profile_requires_roles")
	}
	return nil
}
