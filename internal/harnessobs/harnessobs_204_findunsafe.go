package harnessobs

func findUnsafe(value any) (string, string) {

	return findUnsafeAt("", value)
}
