package main

func hasInstallRepoObserverSubcommand(args []string) bool {
	return len(args) != 0 && args[0] == "repo-observer"
}
