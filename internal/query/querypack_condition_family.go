package query

import "strings"

func familyForForensicCondition(id string) string {
	// familyForForensicCondition keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	switch {
	case strings.Contains(id, "redaction"):
		return EvidenceFamilyRedaction
	case strings.Contains(id, "retention"), strings.Contains(id, "raw_reference"), strings.Contains(id, "critical_evidence"):
		return EvidenceFamilyRetention
	default:
		return EvidenceFamilyClaim
	}
}

func familyForAdapterCondition(id string) string {
	// familyForAdapterCondition keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if strings.Contains(id, "task") {
		return EvidenceFamilyTask
	}
	return nonTaskAdapterFamily(id)
}

func nonTaskAdapterFamily(id string) string {
	// nonTaskAdapterFamily keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	switch {
	case strings.Contains(id, "file"):
		return EvidenceFamilyFileMutations
	case strings.Contains(id, "test"):
		return EvidenceFamilyTest
	default:
		return EvidenceFamilyAdapterCapture
	}
}
