package main

import (
	"strings"
)

func githubAPIURL(apiURLFlag string, getenv func(string) string) (string, error) {
	apiURL := strings.TrimSpace(apiURLFlag)
	if apiURL == "" {
		// GitHub Enterprise runners set GITHUB_API_URL; explicit flags remain
		// higher priority for replay fixtures and local verification.
		apiURL = getenv("GITHUB_API_URL")
	}
	if apiURL == "" {
		apiURL = "https://api.github.com"
	}
	if err := validateGitHubAPIURL(apiURL, getenv("GITHUB_SERVER_URL")); err != nil {
		return "", err
	}
	return strings.TrimRight(apiURL, "/"), nil
}
