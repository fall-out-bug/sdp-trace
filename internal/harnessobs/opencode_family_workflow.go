package harnessobs

func openCodeReviewFamily(signals []string) bool {
	return hasSignal(signals, "review") || hasSignalPrefix(signals, "review.")
}

func openCodePRFamily(signals []string) bool {
	return hasSignal(signals, "pull_request", "pull request") || hasSignalPrefix(signals, "pr.", "pr_")
}

func openCodeMergeFamily(signals []string) bool {
	return hasSignal(signals, "merge") || hasSignalPrefix(signals, "merge.")
}
