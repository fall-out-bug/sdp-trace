package main

import (
	"net/url"
	"strings"
)

func githubAPIHostAllowed(host, serverURL string) bool {
	serverHost := githubServerHost(serverURL)
	if publicGitHubServerHost(serverHost) {
		// Public github.com maps to the public API host; Enterprise hosts must
		// bind exactly to the configured server hostname.
		return host == "api.github.com"
	}
	return host == serverHost
}

func publicGitHubServerHost(serverHost string) bool {
	return serverHost == "" || serverHost == "github.com"
}

func githubServerHost(serverURL string) string {
	if strings.TrimSpace(serverURL) == "" {
		// Empty server URL means public GitHub, handled by githubAPIHostAllowed.
		return ""
	}
	parsed, err := url.Parse(serverURL)
	if err != nil {
		return ""
	}
	return strings.ToLower(parsed.Hostname())
}
