package harnessobs

func openCodeHarnessFamily(signals []string) bool {
	return hasSignal(signals, "session.started", "session.completed", "run.started", "run.completed", "step_start", "step-start", "step_finish", "step-finish") ||
		hasSignalPrefix(signals, "session.", "run.")
}
