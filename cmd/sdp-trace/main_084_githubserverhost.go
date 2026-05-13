package main

import (
	"net/url"
	"strings"
)

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
