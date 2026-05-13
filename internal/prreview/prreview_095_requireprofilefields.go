package prreview

import (
	"errors"

	"strings"
)

func requireProfileFields(profile ReviewProfile) error {
	// requireProfileFields keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
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
