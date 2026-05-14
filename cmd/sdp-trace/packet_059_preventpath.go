package main

func prEventPath(source string, eventPath string, getenv func(string) string) string {
	if source == "github-actions" && eventPath == "" {
		// GitHub Actions events default to the runner-provided event file when
		// the packet command is not using an explicit fixture path.
		return getenv("GITHUB_EVENT_PATH")
	}
	return eventPath
}
