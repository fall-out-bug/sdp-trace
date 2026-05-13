package main

func exportCrossRepoPostureRequested(args []string) bool {
	return exportCommandIs(args, "cross-repo-posture")
}
