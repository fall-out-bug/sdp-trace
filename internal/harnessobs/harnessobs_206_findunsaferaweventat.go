package harnessobs

func findUnsafeRawEventAt(path string, value any) (string, string) {
	return findUnsafeValueAt(path, value, true)
}
