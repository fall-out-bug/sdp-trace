package authority

import (
	"strings"
)

func evidenceRefsReason(refs []string, resolution map[string]string) string {
	// evidenceRefsReason keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if len(refs) == 0 {

		return "evidence_ref_missing"
	}
	for _, ref := range refs {

		if reason := evidenceRefReason(ref, resolution[ref]); reason != "" {
			return reason
		}
	}
	return ""
}

func evidenceRefReason(ref string, resolution string) string {
	// evidenceRefReason keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if malformedEvidenceRef(ref) {

		return "evidence_ref_malformed"
	}
	if reason, ok := evidenceRefResolutionReasons[resolution]; ok {

		return reason
	}
	if unresolvedExternalEvidenceRef(ref, resolution) {
		return "external_evidence_unresolved"
	}
	return ""
}

func malformedEvidenceRef(ref string) bool {
	return unsafeRefPattern.MatchString(ref) || !evidenceRefPattern.MatchString(ref)
}

func unresolvedExternalEvidenceRef(ref string, resolution string) bool {
	return strings.HasPrefix(ref, "external:") && resolution != "resolved"
}
