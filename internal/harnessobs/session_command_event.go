package harnessobs

func sessionCommandModelEvent(session SessionRun) Event {
	return normalizedEvent(
		"session-command-model",
		"model",
		"model_observed",
		sessionCommandObservedAt(session),
		"session-command",
		safeToken(session.CommandModel),
	)
}
