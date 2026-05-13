package harnessobs

func unsafeStringPath(value string, rawEvent bool) bool {
	return !rawEvent && unsafePathValue(value)
}
