package prreview

import (
	"fmt"
)

func validateProfileHeader(profile ReviewProfile) error {
	// validateProfileHeader keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if profile.SchemaVersion != "" && profile.SchemaVersion != SchemaVersionProfile {

		return fmt.Errorf("invalid_profile_schema_version: %s", profile.SchemaVersion)
	}
	return requireProfileFields(profile)
}
