package recorder

func recorderResult(prepared preparedRun, exitCode int) RecorderResult {
	// Results expose only the stable run directory, command exit code, and
	// resolved contract; event details stay in the run artifacts.
	return RecorderResult{
		RunDir:   prepared.runDir,
		ExitCode: exitCode,
		Contract: prepared.contract,
	}
}
