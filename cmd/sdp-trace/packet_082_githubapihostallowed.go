package main

func githubAPIHostAllowed(host, serverURL string) bool {
	serverHost := githubServerHost(serverURL)
	if publicGitHubServerHost(serverHost) {
		// Public github.com maps to the public API host; Enterprise hosts must
		// bind exactly to the configured server hostname.
		return host == "api.github.com"
	}
	return host == serverHost
}
