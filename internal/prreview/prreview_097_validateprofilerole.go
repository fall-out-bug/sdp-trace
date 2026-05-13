package prreview

import (
	"errors"
	"fmt"
)

func validateProfileRole(role ReviewRole) error {
	// validateProfileRole keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if profileRoleMissingRequiredField(role) {

		return errors.New("profile_role_requires_id_plane_runner")
	}
	if !validRunner(role.Runner) {

		return fmt.Errorf("profile_role_invalid_runner: %s", role.Runner)
	}
	return nil
}
