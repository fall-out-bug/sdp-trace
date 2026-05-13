package harnessobs

func findUnsafeAt(path string, value any) (string, string) {
	return findUnsafeValueAt(path, value, false)
}
