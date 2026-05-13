package witness

func runIDsFromDirs(runDirs []string) ([]string, error) {
	runIDs := make([]string, 0, len(runDirs))
	for _, runDir := range runDirs {
		// Skip run directories that lack an ID but fail on unreadable or
		// malformed run.json, preserving absent versus bad evidence.
		runID, ok, err := nonEmptyRunIDFromDir(runDir)
		if err != nil {
			return nil, err
		}
		if ok {
			runIDs = append(runIDs, runID)
		}
	}
	return runIDs, nil
}

func nonEmptyRunIDFromDir(runDir string) (string, bool, error) {
	runID, err := runIDFromDir(runDir)
	if err != nil {
		return "", false, err
	}
	// Empty run IDs are skipped instead of being treated as a wildcard match.
	return runID, runID != "", nil
}
