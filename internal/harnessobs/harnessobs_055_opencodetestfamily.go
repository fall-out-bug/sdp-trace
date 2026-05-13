package harnessobs

func openCodeTestFamily(signals []string) bool {
	return hasSignal(signals, "test.finished", "test.started", "test.passed", "test.failed") ||
		hasSignalPrefix(signals, "test.")
}
