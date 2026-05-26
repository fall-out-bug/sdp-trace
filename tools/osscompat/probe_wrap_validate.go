package main

func validateWrapOutput(stdout []byte, tmpDir string) (string, verifierState, string) {
	runDir, err := parseWrapRunDir(string(stdout))
	if err != nil {
		return "", stateCannotVerify, err.Error()
	}
	if err := validateRunDirUnderTmp(runDir, tmpDir); err != nil {
		return "", stateCannotVerify, err.Error()
	}
	return runDir, "", ""
}
