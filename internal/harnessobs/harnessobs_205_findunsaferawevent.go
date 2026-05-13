package harnessobs

func findUnsafeRawEvent(value any) (string, string) {

	return findUnsafeRawEventAt("", value)
}
