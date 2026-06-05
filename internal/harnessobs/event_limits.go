package harnessobs

// Event limits convert optional profile limits to the package defaults used by
// JSONL scanning.
func effectiveEventLimits(limits Limits) (int, int) {
	return effectiveLineLimit(limits.MaxLineBytes), effectiveEventLimit(limits.MaxEvents)
}

func effectiveLineLimit(maxLine int) int {
	if maxLine <= 0 {
		return DefaultMaxLineBytes
	}
	return maxLine
}

func effectiveEventLimit(maxEvents int) int {
	if maxEvents <= 0 {
		return DefaultMaxEvents
	}
	return maxEvents
}
