package harnessobs

func openCodeHarnessFamily(signals []string) bool {
	return hasSignal(signals, "session.started", "session.completed", "run.started", "run.completed", "step_start", "step-start", "step_finish", "step-finish") ||
		hasSignalPrefix(signals, "session.", "run.")
}

func openCodeInteractionFamily(raw map[string]any, signals []string) bool {
	return hasKey(raw, "role") ||
		hasSignal(signals, "message", "response", "text") ||
		hasSignalPrefix(signals, "message.", "response.")
}
