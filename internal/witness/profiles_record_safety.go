package witness

import (
	"strings"
)

func applyOutputSafetyPass(record *Record) {
	// Preserve the caller's record while refreshing only the output-safety
	// attestation after the serialized trust payload has been scanned.
	if record.OutputSafety == nil {
		record.OutputSafety = &OutputSafety{}
	}
	record.OutputSafety.State = statePass
	record.OutputSafety.VerifiedAbsentClasses = safetyClasses
}

func unsafeOutputRecord(kind string) Record {
	safe := baseRecord(kind)
	// Unsafe output candidates are replaced by a minimal failure record so the
	// published artifact carries the verdict without carrying the unsafe data.
	safe.ProfileStates = defaultProfileStates(stateFail, "cannot_verify")
	safe.OutputSafety = &OutputSafety{State: stateFail, VerifiedAbsentClasses: safetyClasses}
	applyProfileState(&safe, StatusFail, stateFail, ReasonUnsafeOutput)
	return safe
}

func forbiddenOutputPresent(raw []byte) bool {
	// Output safety checks marker families before publication; unsafe content is
	// never echoed in the resulting failure record.
	// Secret-like structured markers are checked before general lowercase marker
	// scanning.
	// Marker matching is deliberately string-based so it works on serialized JSON
	// without schema-specific traversal.
	if containsSecretLike(raw) {
		return true
	}
	// Marker checks catch unsafe payload classes even when they are not shaped
	// like standard tokens or private keys.
	text := strings.ToLower(string(raw))
	for _, marker := range outputSafetyMarkers {
		if strings.Contains(text, marker) {
			return true
		}
	}
	return jwtLike(text)
}
