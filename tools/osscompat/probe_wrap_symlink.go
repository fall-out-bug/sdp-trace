package main

func checkRunJSONUnderTmp(runJSONPath, tmpDir string) error {
	resolvedPath, err := evalSymlinkPath(runJSONPath, "run.json path")
	if err != nil {
		return err
	}
	resolvedTmp, err := evalSymlinkPath(tmpDir, "tmpDir path")
	if err != nil {
		return err
	}
	return requireResolvedPathUnderTmp(resolvedPath, resolvedTmp)
}
