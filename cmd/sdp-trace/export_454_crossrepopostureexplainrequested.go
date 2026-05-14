package main

func exportCrossRepoPostureExplainRequested(args []string) bool {
	return exportCommandIs(args, "cross-repo-posture") && exportSubcommandIs(args, "explain")
}
