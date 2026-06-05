package harnessobs

import "strings"

// safeRef accepts only stable local identifiers or SHA-256 digests for ordinary
// event references. Operation references use a separate prefix-aware rule below.
func safeRef(ref string) bool {
	return safeIDPattern.MatchString(ref) || sha256Pattern.MatchString(ref)
}

// safeOperationRef allows adapter and delivery trace operation namespaces
// without allowing path traversal, URLs, or unbounded reference strings.
func safeOperationRef(ref string) bool {
	if operationRefPrefix(ref) {
		return safePrefixedOperationRef(ref)
	}
	return safeRef(ref)
}

// operationRefPrefix identifies operation namespaces that are not plain safe
// identifiers but are still part of the portable harness observation contract.
func operationRefPrefix(ref string) bool {
	return strings.HasPrefix(ref, "adapter-run:") || strings.HasPrefix(ref, "delivery-trace:")
}

// safePrefixedOperationRef keeps prefixed operation references textual and
// bounded so they cannot smuggle local paths or external URLs into evidence.
func safePrefixedOperationRef(ref string) bool {
	return !strings.Contains(ref, "..") && !strings.Contains(ref, "://") && len(ref) <= 256
}

// safeEvent renders invalid event IDs as a stable placeholder in validation
// diagnostics instead of echoing potentially unsafe source content.
func safeEvent(eventID string) string {
	if safeIDPattern.MatchString(eventID) {
		return eventID
	}

	return "unknown_event"
}
