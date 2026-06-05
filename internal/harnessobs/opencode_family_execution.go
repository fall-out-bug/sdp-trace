package harnessobs

func openCodeTestFamily(signals []string) bool {
	return hasSignal(signals, "test.finished", "test.started", "test.passed", "test.failed") ||
		hasSignalPrefix(signals, "test.")
}

func openCodePhaseFamily(raw map[string]any, signals []string) bool {
	return hasKey(raw, "phase") ||
		hasSignal(signals, "phase", "gsd.phase_path", "gsd-plan-phase") ||
		hasSignalPrefix(signals, "phase.", "gsd.", "gsd_")
}
