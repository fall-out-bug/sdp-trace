package harnessobs

func writeSessionJSON(path string, run SessionRun) error {
	return writeJSON(path, run)
}
