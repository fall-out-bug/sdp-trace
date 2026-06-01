package harnessobs

func writeJSON(path string, value any) error {
	data, err := jsonArtifactData(value)
	if err != nil {
		return err
	}
	return writeJSONFile(path, data)
}
