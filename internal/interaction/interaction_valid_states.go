package interaction

func validEventType(value string) bool {
	return frictionClass(value) != ""
}

func validRetention(value string) bool {
	// validRetention keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	switch value {
	case RetentionDigestOnly, RetentionSanitizedExcerpt, RetentionEncryptedRawRef, RetentionExternalArtifactRef, RetentionNotAssessed:
		return true
	default:
		return false
	}
}

func validCompleteness(value string) bool {
	// validCompleteness keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	switch value {
	case CompletenessComplete, CompletenessPartial, CompletenessNotAssessed, CompletenessCannotVerify:
		return true
	default:
		return false
	}
}

func validChannelExclusivity(value string) bool {
	// validChannelExclusivity keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	switch value {
	case ChannelExclusivityNotAssessed, StateReferenced, StateCannotVerify:
		return true
	default:
		return false
	}
}
func validEventState(value string) bool {
	// validEventState keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	switch value {
	case StateObserved, StateReferenced, StateUnreferenced, StateNotAssessed, StateCannotVerify, "redacted":
		return true
	default:
		return false
	}
}

func frictionClass(eventType string) string {
	return frictionClasses[eventType]
}
