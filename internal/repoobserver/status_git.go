package repoobserver

func statusGitRefs(repoRoot string) (string, string) {
	// Git refs are local structural observations; empty strings mean git could
	// not provide that field during this snapshot.
	return gitOutput(repoRoot, "rev-parse", "--verify", "HEAD"),
		gitOutput(repoRoot, "branch", "--show-current")
}
