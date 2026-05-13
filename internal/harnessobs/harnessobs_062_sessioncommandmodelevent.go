package harnessobs

func sessionCommandModelEvent(session SessionRun) Event {
	// sessionCommandModelEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return normalizedEvent(
		"session-command-model",
		"model",
		"model_observed",
		sessionCommandObservedAt(session),
		"session-command",
		safeToken(session.CommandModel),
	)
}
