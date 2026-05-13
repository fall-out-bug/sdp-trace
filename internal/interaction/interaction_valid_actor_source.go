package interaction

func validActorType(value string) bool {
	// validActorType keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	switch value {
	case "human_user", "human_role", "human_group", "ai_agent", "model", "system", "tool", "organization", "other":
		return true
	default:
		return false
	}
}

func validSourceType(value string) bool {
	// validSourceType keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	switch value {
	case SourceObservedControlChannel, SourcePreclassifiedTranscript:
		return true
	default:
		return false
	}
}
