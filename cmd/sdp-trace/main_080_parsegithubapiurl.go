package main

import (
	"fmt"
	"net/url"
)

func parseGitHubAPIURL(apiURL string) (*url.URL, error) {
	parsed, err := url.Parse(apiURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("unsafe GitHub API URL %q", apiURL)
	}
	// The parsed URL is still untrusted until scheme, credentials, and host are
	// validated by the caller.
	return parsed, nil
}
