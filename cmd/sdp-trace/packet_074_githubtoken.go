package main

func githubToken(getenv func(string) string) string {
	token := getenv("GITHUB_TOKEN")
	if token == "" {
		// GH_TOKEN is a local CLI fallback only; callers still validate that a
		// token exists before making artifact API requests.
		return getenv("GH_TOKEN")
	}
	return token
}
