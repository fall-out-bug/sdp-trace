package repoobserver

type targetFile struct {
	path       string
	content    string
	executable bool
}

func installTargets(opts Options, repoID string) []targetFile {
	// Generated files stay portable: hooks, local config/docs, and a GitHub
	// workflow declaration only.
	return []targetFile{
		{path: ".sdp-trace/README.md", content: sdpTraceReadme()},
		{path: ".sdp-trace/config.json", content: sdpTraceConfig(opts, repoID)},
		{path: ".githooks/pre-commit", content: hookScript("pre-commit"), executable: true},
		{path: ".githooks/post-commit", content: hookScript("post-commit"), executable: true},
		{path: ".githooks/pre-push", content: hookScript("pre-push"), executable: true},
		{path: ".github/workflows/sdp-trace-observe.yml", content: githubWorkflow()},
	}
}
